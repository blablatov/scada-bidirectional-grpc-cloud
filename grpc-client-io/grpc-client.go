// gRPC-клиент
// go build . && go run .

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	pb "github.com/blablatov/scada-bidirectional-grpc-cloud/grpc-cloud-proto"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/encoding/gzip"
)

var (
	crtFile = filepath.Join("..", "grpc-certs", "client.crt")
	keyFile = filepath.Join("..", "grpc-certs", "client.key")
	caFile  = filepath.Join("..", "grpc-certs", "ca.crt")
	// Создаем реестр метрик и стандартные клиентские метрики
	reg               = prometheus.NewRegistry()
	grpcMetrics       = grpc_prometheus.NewClientMetrics()
	custMetricCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cloud_client_Id_count",
		Help: "Total number of RPCs Id on the client.",
	}, []string{"Id"})
)

const (
	address  = "localhost:50051"
	hostname = "localhost"
)

// Регистрируем стандартные клиентские метрики и добавленный сборщик в реестре
func init() {
	reg.MustRegister(grpcMetrics, custMetricCounter)
}

func main() {
	log.SetPrefix("Client event: ")
	log.SetFlags(log.Lshortfile)

	// Set up the credentials for the connection
	// Значение токена OAuth2. Используем строку, прописанную в коде
	autok := oauth.NewOauthAccess(fetchToken())

	// Load the client certificates from disk
	// Создаем пары ключей X.509 непосредственно из ключа и сертификата сервера
	certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("could not load client key pair: %s", err)
	}

	// Create a certificate pool from the certificate authority
	// Генерируем пул сертификатов в нашем локальном удостоверяющем центре
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	// Append the certificates from the CA
	// Добавляем клиентские сертификаты из локального удостоверяющего центра в сгенерированный пул
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append ca certs")
	}

	opts := []grpc.DialOption{
		// Register interceptor of stream. Регистрация потокового перехватчика
		grpc.WithStreamInterceptor(clientStreamInterceptor),

		// Перехватчик метрик для потокового клиента. Stream interceptor of metrics
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),

		// Указываем один и тот же токен OAuth в параметрах всех вызовов в рамках одного соединения
		// Если нужно указывать токен для каждого вызова отдельно, используем CallOption
		grpc.WithPerRPCCredentials(autok),

		// Transport credentials
		// Указываем транспортные аутентификационные данные в виде параметров соединения
		// Поле ServerName должно быть равно значению Common Name, указанному в сертификате
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:   hostname,
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})),
	}

	// Set up a connection to the server
	// Устанавливаем безопасное соединение с сервером, передаем параметры аутентификации
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Создаем HTTP-сервер для Prometheus.
	// HTTP-путь для сбора метрик начинается с /metrics и находится на порту 9094
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9094)}

	// Запускаем на клиенте HTTP-сервер для Prometheus
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	// Передаем соединение и создаем заглушку.
	// Ее экземпляр содержит все удаленные методы, которые можно вызвать на сервере
	client := pb.NewCloudExchangeClient(conn)

	// Finding of Duration. Тестированием определить оптимальное значение для крайнего срока кпд
	clientDeadline := time.Now().Add(time.Duration(5000 * time.Millisecond))
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	// Process Order : Bi-distreaming scenario
	// Вызываем удаленный метод и получаем ссылку на поток записи и чтения на клиентской стороне
	streamCloud, err := client.ProcessCloud(ctx, grpc.UseCompressor(gzip.Name))
	if err != nil {
		log.Fatalf("%v.ProcessCloud(_) = _, %v", client, err)
	}

	mp := map[string]float64{
		"102": 42.4,
		"103": 53.1,
		"104": 24.7,
		"105": 35.2,
		"11":  11.9,
	}

	for k, v := range mp {
		if k != "" && k != "0" {
			// Sends IDs. Отправляем сообщения с ID сервису
			if err := streamCloud.Send(&pb.RequestIO{Id: k, Measurement: v}); err != nil {
				log.Fatalf("%v.Send(%v) = %v", client, k, err)
			}
		} else {
			log.Printf("ID not found(%s) = %b", k, v)
		}
	}

	chs := make(chan struct{}) // Создаем канал для горутин (create chanel for goroutines)

	// Вызываем функцию с помощью горутин, распараллеливаем чтение сообщений, возвращаемых сервисом
	go func() {
		asncClientBidirectionalRPC(streamCloud, chs)
		chs <- struct{}{}
		close(chs)
	}()

	time.Sleep(time.Millisecond * 1000) //  Wait time. Имитируем задержку при отправке сервису сообщений.

	// Сигнализируем о завершении клиентского потока (с ID заказов)
	// Signal about close stream of client
	if err := streamCloud.CloseSend(); err != nil {
		log.Fatal(err)
	}

	// Cancelling the RPC. Отмена удаленного вызова gRPC на клиентской стороне
	cancel()
	log.Printf("RPC Status : %v", ctx.Err()) // Status of context. Состояние текущего контекста

	<-chs
}

func asncClientBidirectionalRPC(streamCloud pb.CloudExchange_ProcessCloudClient, c chan struct{}) {
	for {
		// Read messages on side of client
		// Читаем сообщения сервиса на клиентской стороне
		combStatus, errCloud := streamCloud.Recv()

		if errCloud != nil {
			log.Printf("Error Receiving messages: %v", errCloud)
			break
		} else {
			if errCloud == io.EOF { // End of stream. Обнаружение конца потока.
				break
			}
			log.Println("Combined status: ", combStatus.Status, combStatus.IOList)
		}
	}
	<-c
}

// Provides OAuth2 connection token
// Учетные данные для соединения. Предоставление токена OAuth2
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "blablatok-tokblabla-blablatok",
	}
}

// Client stream interceptor in gRPC
// Клиентский потоковый перехватчик в gRPC
func clientStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
	method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

	// Preprocessing stage, haves access to RPC request before sent to server
	// Этап предобработки, есть доступ к RPC-запросу перед его отправкой на сервер
	log.Println("===== [Client Interceptor] ", method)

	s, err := streamer(ctx, desc, cc, method, opts...) // Call func streamer. Вызов функции streamer
	if err != nil {
		return nil, err
	}
	// Creating wrapper around Client Stream interface, with intercept and go back to app
	// Создание обертки вокруг интерфейса ClientStream, с перехватом и возвращением приложению
	return newWrappedStream(s), nil
}

// Wrapper for interface of rpc.ClientStream
// Обертка для интерфейса grpc.ClientStream
type wrappedStream struct {
	grpc.ClientStream
}

// Func for intercepting received messages of streaming gRPC
// Функция для перехвата принимаемых сообщений потокового gRPC
func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("===== [Client Stream Interceptor] Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

// Func for intercepting sended messages of streaming gRPC
// Функция для перехвата отправляемых сообщений потокового gRPC
func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("===== [Client Stream Interceptor] Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}
