// gRPC-сервис
// go build . && go run .

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/blablatov/scada-bidirectional-grpc-cloud/grpc-cloud-proto"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	crtFile             = filepath.Join("..", "grpc-certs", "server.crt")
	keyFile             = filepath.Join("..", "grpc-certs", "server.key")
	caFile              = filepath.Join("..", "grpc-certs", "ca.crt")
	errMissingMetadata  = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken     = status.Errorf(codes.Unauthenticated, "invalid token")
	reg                 = prometheus.NewRegistry()
	grpcMetrics         = grpc_prometheus.NewServerMetrics()
	customMetricCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cloud_server_Id_count",
		Help: "Total number of RPCs Id on the server.",
	}, []string{"Id"})
)

const (
	port         = ":50051"
	reqBatchSize = 1 // Group of orders. Заказы обрабатываются группами.
)

// Регистрируем стандартные серверные метрики и добавленный сборщик в реестре
func init() {
	reg.MustRegister(grpcMetrics, customMetricCounter)
}

func main() {
	log.SetPrefix("Server event: ")
	log.SetFlags(log.Lshortfile)

	// Reading opened/closed keys to enable TLS
	// Считываем и анализируем открытый/закрытый ключи и создаем сертификат, чтобы включить TLS
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	// Create a certificate pool from the certificate authority
	// Генерируем пул сертификатов в удостоверяющем центре
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	// Append the client certificates from the CA
	// Добавляем клиентские сертификаты из удостоверяющего центра в сгенерированный пул
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append client certs")
	}

	opts := []grpc.ServerOption{
		// Enable TLS for all incoming connections. Включаем TLS для всех входящих соединений путем
		grpc.Creds( // Create the TLS credentials. Создание аутентификационных данных TLS
			credentials.NewTLS(&tls.Config{
				ClientAuth:   tls.RequireAndVerifyClientCert,
				Certificates: []tls.Certificate{cert},
				ClientCAs:    certPool,
			},
			)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// Registers unary interceptor to gRPC-server. Регистрация унарного перехватчика
			// Будет направлять клиентские запросы к функции ensureValidBasicCredentials
			grpc.UnaryServerInterceptor(ensureValidToken),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// Регистрация дополнительного потокового перехватчика на gRPC-сервере
			// Направляет клиентские запросы к функции orderServerStreamInterceptor
			grpc.StreamServerInterceptor(orderServerStreamInterceptor),

			// Перехватчик метрик для потокового сервера. Stream interceptor of metrics
			grpcMetrics.StreamServerInterceptor(),
		)),
	}

	// // Creates new gRPC server, sends him data of authentification
	// // Создаем новый экземпляр gRPC-сервера, передавая ему аутентификационные данные
	grpcServer := grpc.NewServer(opts...)

	initIOData()

	// Создаем HTTP-сервер для Prometheus.
	// HTTP-путь для сбора метрик начинается с /metrics и находится на порту 9092
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9092)}

	// Инициализируем все стандартные метрики
	pb.RegisterCloudExchangeServer(grpcServer, &mserver{})
	grpcMetrics.InitializeMetrics(grpcServer)

	// Запускаем на сервере HTTP-сервер для Prometheus
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	// Начинаем прослушивать TCP на порту 50051. Listen on TCP port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Starting gRPC listener on port " + port)

	// Binds gRPC server to listener, waiting for messages on port 50051
	// Привязываем gRPC-сервер к прослушивателю, ожидающему сообщений на порту 50051
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initIOData() {
	reqMap["102"] = pb.RequestIO{Id: "102", Sensors: []string{"Dallas", "Texas Instruments"}, Description: "sensor#99", Destination: "Surgut, City", Measurement: 55}
	reqMap["103"] = pb.RequestIO{Id: "103", Sensors: []string{"Dallas semiconductor", "Texas Instruments"}, Description: "sensor#77", Destination: "Samara, City", Measurement: 22}
	reqMap["104"] = pb.RequestIO{Id: "104", Sensors: []string{"Dallas", "Texas Instruments"}, Description: "sensor#99", Destination: "Moscow, City", Measurement: 35}
	reqMap["105"] = pb.RequestIO{Id: "105", Sensors: []string{"Dallas semiconductor", "Inc", "Texas Instruments"}, Description: "sensor#33", Destination: "Volgograd, City", Measurement: 38}

}

// Validates the authorization. Валидация токена
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Performs validation of token matching an arbitrary string
	// Выполняем проверку токена, соответствующего нашему произвольно заданному
	return token == "blablatok-tokblabla-blablatok"
}

// Определяем функцию ensureValidToken для проверки подлинности токена. Validate token
// Если тот отсутствует или недействителен, тогда перехватчик блокирует выполнение и возвращает ошибку
// Или вызывается следующий обработчик, которому передается контекст и интерфейс
func ensureValidToken(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token
	// Если токен действителен, продолжаем выполнение обработчика
	return handler(ctx, req)
}

// Stream wraps around the embedded grpc.Server Stream, and intercepts the RecvMsg and SendMsg method call
// Обертка вокруг встроенного интерфейса grpc.ServerStream, перехватывает вызовы методов RecvMsg и SendMsg
type wrappedStream struct {
	grpc.ServerStream
}

// RecvMsg wrapper function, handles received gRPC streaming messages
// Функция обертки RecvMsg, обрабатывает принимаемые сообщения потокового gRPC
func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("====== [Server Stream Interceptor Wrapper] Receive a message (Body: %s) at %s",
		m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

// The wrapper function SendMsg, handles the sent messages of the streaming gRPC
// Функция обертки SendMsg, обрабатывает отправляемые сообщения потокового gRPC
func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("====== [Server Stream Interceptor Wrapper] Send a message (Body: %s) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

// Creating wrapper function. Создание экземпляра функции-обертки
func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// Creates streaming interceptor. Реализация потокового перехватчика
func orderServerStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	// Pre-processing. Этап предобработки
	log.Println("====== [Server Stream Interceptor] ", info.FullMethod)

	// Invoking the StreamHandler to complete the execution of RPC invocation
	// Вызов метода потокового RPC с помощью обертки.
	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return err
}
