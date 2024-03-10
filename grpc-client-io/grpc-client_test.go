// This conventional test. Традиционный тест клиента.
// Перед выполнением запустить grpc-сервер
// Before his execute run grpc-server
// go run . && go test .

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	pb "github.com/blablatov/scada-bidirectional-grpc-cloud/grpc-cloud-proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/encoding/gzip"
)

// Conventional test that starts a gRPC client test the service with RPC
// Традиционный тест, который запускает клиент для проверки удаленного метода сервиса
func TestProcessCloud(t *testing.T) {
	log.SetPrefix("Client-test event: ")
	log.SetFlags(log.Lshortfile)

	tokau := oauth.NewOauthAccess(fetchToken())

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

	// Указываем аутентификационные данные для транспортного протокола с помощью DialOption
	opts := []grpc.DialOption{
		// Указываем один и тот же токен OAuth в параметрах всех вызовов в рамках одного соединения
		// Если нужно указывать токен для каждого вызова отдельно, используем CallOption
		grpc.WithPerRPCCredentials(tokau),

		// Перехватчик метрик для потокового клиента. Stream interceptor of metrics
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),

		// Указываем транспортные аутентификационные данные в виде параметров соединения
		// Поле ServerName должно быть равно значению Common Name, указанному в сертификате
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:   hostname, // NOTE: this is required!
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewCloudExchangeClient(conn)

	// Finding of Duration. Тестированием определить оптимальное значение для крайнего срока кпд
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	clientDeadline := time.Now().Add(time.Duration(5000 * time.Millisecond))
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

	defer cancel()

	// Process Order : Bi-distreaming scenario
	// Вызываем удаленный метод и получаем ссылку на поток записи и чтения на клиентской стороне
	streamCloud, err := client.ProcessCloud(ctx, grpc.UseCompressor(gzip.Name))
	if err != nil {
		log.Fatalf("%v.ProcessCloud(_) = _, %v", client, err)
	}

	// IDs for test. Мапа с тестируемыми ID
	mp := map[string]float64{
		"102": 42.4,
		"103": 53.1,
		"104": 24.7,
		"105": 35.2,
		"11":  11.9,
	}

	for k, v := range mp {
		if k != "" && v != 0 {
			// Отправляем сообщения сервису
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

	time.Sleep(time.Millisecond * 1000) // Имитируем задержку при отправке сервису сообщений. Wait time

	// Сигнализируем о завершении клиентского потока (с ID заказов)
	// Signal about close stream of client
	if err := streamCloud.CloseSend(); err != nil {
		log.Fatal(err)
	}

	chs <- struct{}{}
}

// Тестирование производительности в цикле за указанное колличество итераций
func BenchmarkProcessCloud(b *testing.B) {

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

	b.ReportAllocs()
	for i := 0; i < 550; i++ {
		tokau := oauth.NewOauthAccess(fetchToken())

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
			grpc.WithPerRPCCredentials(tokau),
			// Указываем транспортные аутентификационные данные в виде параметров соединения
			// Поле ServerName должно быть равно значению Common Name, указанному в сертификате
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				ServerName:   hostname, // NOTE: this is required!
				Certificates: []tls.Certificate{certificate},
				RootCAs:      certPool,
			})),
		}

		conn, err := grpc.Dial(address, opts...) // Подключаемся к серверному приложению
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewCloudExchangeClient(conn)

		// Finding of Duration. Тестированием определить оптимальное значение для крайнего срока кпд
		clientDeadline := time.Now().Add(time.Duration(2000 * time.Millisecond))
		ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

		defer cancel()

		// Process Order : Bi-distreaming scenario
		// Вызываем удаленный метод и получаем ссылку на поток записи и чтения на клиентской стороне
		streamCloud, err := client.ProcessCloud(ctx, grpc.UseCompressor(gzip.Name))
		if err != nil {
			log.Fatalf("%v.ProcessCloud(_) = _, %v", client, err)
		}

		// IDs for test. Мапа с тестируемыми ID
		mp := map[string]float64{
			"102": 42.4,
			"103": 53.1,
			"104": 24.7,
			"105": 35.2,
			"11":  11.9,
		}

		for k, v := range mp {
			if k != "" && v != 0 {
				// Отправляем сообщения сервису
				if err := streamCloud.Send(&pb.RequestIO{Id: k, Measurement: v}); err != nil {
					log.Fatalf("%v.Send(%v) = %v", client, k, err)
				}
			} else {
				log.Printf("ID not found(%s) = %b", k, v)
			}
		}

		chs := make(chan struct{}) // Создаем канал для горутин (create chanel for goroutines)
		// Вызываем функцию с помощью горутин, распараллеливаем чтение сообщений, возвращаемых сервисом
		go asncClientBidirectionalRPC(streamCloud, chs)
		time.Sleep(time.Millisecond * 1000) // Имитируем задержку при отправке сервису сообщений. Wait time

		// Сигнализируем о завершении клиентского потока (с ID заказов)
		// Signal about close stream of client
		if err := streamCloud.CloseSend(); err != nil {
			log.Fatal(err)
		}

		// Cancelling the RPC. Отмена удаленного вызова gRPC на клиентской стороне
		//cancel()
		log.Printf("RPC Status : %v", ctx.Err()) // Status of context. Состояние текущего контекста

		chs <- struct{}{}
	}
}

// Provides OAuth2 connection token
// Тест предоставления токена OAuth2
func TestFetchToken(t *testing.T) {
	go func() {
		fetchToken()
		var tok oauth2.Token
		if tok.AccessToken == "blablatok-tokblabla-blablatok" {
			log.Println("Token true")
		}
	}()
}

// Сигнатура для тестирования типа Streamer. Signature of method to test
type strmer interface {
	Streamer(context.Context, *grpc.StreamDesc, *grpc.ClientConn, string, ...grpc.CallOption) (grpc.ClientStream, error)
}

type strm struct {
	strmer
	ctx     context.Context
	desc    *grpc.StreamDesc
	cc      *grpc.ClientConn
	method  string
	opts    []grpc.CallOption
	success bool
}

func (d *strm) Streamer(context.Context, *grpc.StreamDesc, *grpc.ClientConn, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	if d.success {
		return nil, nil
	}
	return nil, fmt.Errorf("Streamer test error")
}

func TestClientStreamInterceptor(t *testing.T) {

	var su strm
	log.Println("===== [Client Interceptor] ", su.method)
	// Call func streamer. Вызов сигнатуры streamer
	_, err := su.Streamer(su.ctx, su.desc, su.cc, su.method, su.opts...)
	if err != nil {
		log.Println("Error Interceptor", err)
	}
}
