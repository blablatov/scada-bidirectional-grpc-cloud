package mock_grpc_cloud_proto

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	pb "github.com/blablatov/scada-bidirectional-grpc-cloud/grpc-cloud-proto"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
)

// Реализация интерфейса сопоставления.
// rpcMsg implements the gomock.Matcher interface
type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestClient_ProcessCloud(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockCloudExchangeClient(ctrl)

	reqMap["102"] = pb.RequestIO{Id: "102", Sensors: []string{"Dallas", "Texas Instruments"}, Description: "sensor#99", Destination: "Surgut, City", Measurement: 55}
	reqMap["103"] = pb.RequestIO{Id: "103", Sensors: []string{"Dallas semiconductor", "Texas Instruments"}, Description: "sensor#77", Destination: "Samara, City", Measurement: 22}
	reqMap["104"] = pb.RequestIO{Id: "104", Sensors: []string{"Dallas", "Texas Instruments"}, Description: "sensor#99", Destination: "Moscow, City", Measurement: 35}
	reqMap["105"] = pb.RequestIO{Id: "105", Sensors: []string{"Dallas semiconductor", "Inc", "Texas Instruments"}, Description: "sensor#33", Destination: "Volgograd, City", Measurement: 38}

	mockClient.EXPECT().ProcessCloud(gomock.Any(), &rpcMsg{msg: req}).
		Return(&wrappers.StringValue{Value: "Product_mock" + id}, nil)
	testClient_ProcessCloud(t, mockClient)
}

func testClient_ProcessCloud(t *testing.T, client pb.CloudExchangeClient) {

	// Finding of Duration. Тестированием определить оптимальное значение для крайнего срока кпд
	clientDeadline := time.Now().Add(time.Duration(5000 * time.Millisecond))
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

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

	channel := make(chan int) // Создаем канал для горутин (goroutines).
	// Вызываем функцию с помощью горутин, чтобы распараллелить чтение сообщений, возвращаемых сервисом.
	//go asncClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000) // Имитируем задержку при отправке сервису некоторых сообщений

	if err := streamCloud.CloseSend(); err != nil { // Сигнализируем о завершении клиентского потока (с ID заказов).
		log.Fatal(err)
	}
	channel <- 1
}
