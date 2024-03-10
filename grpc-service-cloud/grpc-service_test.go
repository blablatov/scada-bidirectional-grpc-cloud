// Testing remote functions without using network
// Модульное тестирование бизнес-логики удаленных методов без передачи по сети.
// С запуском стандартного gRPC-сервера поверх HTTP/2 на реальном порту.
// Имитация запуска сервера с использованием буфера.

package main

import (
	"context"
	"io"
	"log"
	"net"
	"testing"
	"time"

	pb "github.com/blablatov/scada-bidirectional-grpc-cloud/grpc-cloud-proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"
)

const (
	address = "localhost:50051"
	bufSize = 1024 * 1024
)

var listener *bufconn.Listener

func initGRPCServerHTTP2() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCloudExchangeServer(s, &mserver{})
	//initIOData()
	// Register reflection service on gRPC server.
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

// Initialization of BufConn. Package bufconn provides a net
// Conn implemented by a buffer and related dialing and listening functionality
// Реализует имитацию запуска сервера на реальном порту с использованием буфера
func initGRPCServerBuffConn() {
	listener = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterCloudExchangeServer(s, &mserver{})
	//initIOData()
	// Register reflection service on gRPC server.
	reflection.Register(s)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

// Conventional test that starts a gRPC server and client test the service with RPC
func TestProcessOrders(t *testing.T) {

	log.SetPrefix("Client-test event: ")
	log.SetFlags(log.Lshortfile)

	// Starting a conventional gRPC server runs on HTTP2
	// Запускаем стандартный gRPC-сервер поверх HTTP/2
	initGRPCServerHTTP2()
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // Подключаемся к серверному приложению
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Ее экземпляр содержит все удаленные методы, которые можно вызвать на сервере
	client := pb.NewCloudExchangeClient(conn)

	// Finding of Duration. Тестированием определить оптимальное значение для крайнего срока кпд
	clientDeadline := time.Now().Add(time.Duration(10000 * time.Millisecond))
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

	defer cancel()

	// Process Cloud : Bi-distreaming scenario
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
			// Иммитируем отправку сообщений сервису.
			if err := streamCloud.Send(&pb.RequestIO{Id: k, Measurement: v}); err != nil {
				log.Fatalf("%v.Send(%v) = %v", client, v, err)
			}
		} else {
			log.Printf("ID not found(%s) = %b", k, v)
		}
	}

	channel := make(chan int) // Создаем канал для горутин (create chanel for goroutines)
	// Вызываем функцию с помощью горутин, распараллеливаем чтение сообщений, возвращаемых сервисом
	go asncClientBidirectionalRPC(streamCloud, channel)
	time.Sleep(time.Millisecond * 500) // Имитируем задержку при отправке сервису сообщений. Wait time

	// Сигнализируем о завершении клиентского потока (с ID заказов)
	// Signal about close stream of client
	if err := streamCloud.CloseSend(); err != nil {
		log.Fatal(err)
	}
	channel <- 1
}

// Test written using Buffconn
func TestProcessCloudBufConn(t *testing.T) {
	ctx := context.Background()
	initGRPCServerBuffConn()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Передаем соединение и создаем заглушку
	// Ее экземпляр содержит все удаленные методы, которые можно вызвать на сервере
	client := pb.NewCloudExchangeClient(conn)

	// Finding of Duration. Тестированием определить оптимальное значение для крайнего срока кпд
	clientDeadline := time.Now().Add(time.Duration(10000 * time.Millisecond))
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

	defer cancel()

	// Process Cloud : Bi-distreaming scenario
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
			// Иммитируем отправку сообщений сервису.
			if err := streamCloud.Send(&pb.RequestIO{Id: k, Measurement: v}); err != nil {
				log.Fatalf("%v.Send(%v) = %v", client, v, err)
			}
		} else {
			log.Printf("ID not found(%s) = %b", k, v)
		}
	}

	channel := make(chan int) // Создаем канал для горутин (create chanel for goroutines)
	// Вызываем функцию с помощью горутин, распараллеливаем чтение сообщений, возвращаемых сервисом
	go asncClientBidirectionalRPC(streamCloud, channel)
	time.Sleep(time.Millisecond * 1000) // Имитируем задержку при отправке сервису сообщений. Wait time

	// Сигнализируем о завершении клиентского потока (с ID заказов)
	// Signal about close stream of client
	if err := streamCloud.CloseSend(); err != nil {
		log.Fatal(err)
	}
	channel <- 1
}

// Benchmark test
// Тестирование производительности в цикле за указанное колличество итераций
func BenchmarkProcessOrdersBufConn(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 150; i++ {
		ctx := context.Background()
		initGRPCServerBuffConn()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		// Ее экземпляр содержит все удаленные методы, которые можно вызвать на сервере
		client := pb.NewCloudExchangeClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// Process Cloud : Bi-distreaming scenario
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
				// Отправляем сообщения сервису.
				if err := streamCloud.Send(&pb.RequestIO{Id: k, Measurement: v}); err != nil {
					log.Fatalf("%v.Send(%v) = %v", client, v, err)
				}
			} else {
				log.Printf("ID not found(%s) = %b", k, v)
			}
		}

		channel := make(chan int) // Создаем канал для горутин (create chanel for goroutines)
		// Вызываем функцию с помощью горутин, распараллеливаем чтение сообщений, возвращаемых сервисом
		go asncClientBidirectionalRPC(streamCloud, channel)
		time.Sleep(time.Millisecond * 1000) // Имитируем задержку при отправке сервису сообщений. Wait time

		// Сигнализируем о завершении клиентского потока (с ID заказов)
		// Signal about close stream of client
		if err := streamCloud.CloseSend(); err != nil {
			log.Fatal(err)
		}
		channel <- 1
	}
}

func asncClientBidirectionalRPC(streamCloud pb.CloudExchange_ProcessCloudClient, c chan int) {
	for {
		// Читаем сообщения сервиса на клиентской стороне
		// Read messages on side of client
		combStatus, errCloud := streamCloud.Recv()

		if errCloud != nil {
			log.Printf("Error Receiving messages: %v", errCloud)
			break
		} else {
			if errCloud == io.EOF { // Обнаружение конца потока. End of stream
				break
			}
			log.Println("Combined status : ", combStatus.IOList)
		}
	}
	<-c
}
