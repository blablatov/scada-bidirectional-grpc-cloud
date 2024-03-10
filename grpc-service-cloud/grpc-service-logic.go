// Методы (логика) gRPC-сервиса

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	//sl "github.com/blablatov/scada-bidirectional-grpc-cloud/clickhouse/selectlog"
	in "github.com/blablatov/scada-bidirectional-grpc-cloud/clickhouse/insertlog"
	pb "github.com/blablatov/scada-bidirectional-grpc-cloud/grpc-cloud-proto"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var reqMap = make(map[string]pb.RequestIO)

// mСервер реализует order_management
type mserver struct {
	reqMap map[string]*pb.RequestIO
}

// Bi-directional Streaming RPC
// Двунаправленный потоковый RPC
func (s *mserver) ProcessCloud(stream pb.CloudExchange_ProcessCloudServer) error {

	batchMarker := 1
	var statusMap = make(map[string]pb.StatusIO)
	for {

		switch {
		case stream.Context().Err() == context.Canceled:
			log.Printf("Context Cacelled for this stream: -> %s", stream.Context().Err())
			log.Printf("Stopped processing any more order of this stream!")
			return stream.Context().Err()

		case stream.Context().Err() == context.DeadlineExceeded:
			log.Printf("Deadline was exceeded for this stream: -> %s", stream.Context().Err())
			log.Printf("Stopped processing any more order of this stream!")
			return stream.Context().Err()

		default:
			// Err of ID. Проверка ID
			reqId, err := stream.Recv() // Reads data from the in stream. Читает данные из входящего потока
			log.Printf("Reading Proc order : %s", reqId)

			// Добавление новой метрики с помощью созданного счетчика
			customMetricCounter.WithLabelValues(reqId.Id).Inc()
			customMetricCounter.WithLabelValues(reqId.Sensor).Inc()
			customMetricCounter.WithLabelValues(reqId.Description).Inc()
			customMetricCounter.WithLabelValues(reqId.Destination).Inc()

			if reqId == nil {
				return err
			}

			if k, ok := reqMap[reqId.Id]; !ok {
				log.Printf("Request ID is invalid! -> Received Request ID %v", k)
				errStatus := status.New(codes.InvalidArgument, "Request ID received is not found - Invalid information")
				ds, err := errStatus.WithDetails(
					&epb.BadRequest_FieldViolation{
						Field:       "ID",
						Description: fmt.Sprintf("Request ID received is not found %s : %s", reqId, reqId.Id),
					},
				)
				if err == nil {
					return errStatus.Err()
				}
				return ds.Err()
			}

			if reqId.String() == `value:"-1"` {
				log.Printf("Request ID is invalid! -> Received Request ID %s", reqId)

				errStatus := status.New(codes.InvalidArgument, "Request ID received is not valid  - Invalid information")
				ds, err := errStatus.WithDetails(
					&epb.BadRequest_FieldViolation{
						Field:       "ID",
						Description: fmt.Sprintf("Request ID received is not valid %s : %s", reqId, reqId.Id),
					},
				)
				if err == nil {
					return errStatus.Err()
				}
				return ds.Err()
			}

			// Checks to Err EOF
			if err == io.EOF { // Reads IDs to EOF. Продолжаем читать, пока не обнаружим конец потока
				// Client has sent all the messages. Send remaining shipments
				log.Printf("EOF : %s", reqId)
				for _, shipment := range statusMap {
					// If EOF sends all data of groups
					// При обнаружении конца потока отправляем клиенту все сгруппированные оставшиеся данные
					if err := stream.Send(&shipment); err != nil {
						return err
					}
				}
				return nil //Closes stream. Сервер завершает поток, возвращая nil
			}
			if err != nil {
				log.Println(err)
				return err
			}

			// Вызов метода выполнения запроса insert к СУБД ClickHouse
			chs := make(chan string, 1)
			chm := make(chan float64, 1)
			chi := make(chan bool, 1)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() error {
				v := new(in.InsData)

				v.Mu.Lock()
				chs <- reqId.Id
				v.Mu.Unlock()

				v.Mu.Lock()
				chm <- reqId.Measurement
				v.Mu.Unlock()

				err := v.InsertLog(wg, chs, chm, chi)
				if err != nil {
					log.Printf("Error connect to dBase: %v", err)
					return err
				}
				log.Printf("Id inserted[%s]: %v", <-chs, <-chi)
				return nil
			}()
			go func() {
				wg.Wait()
				close(chs)
				close(chi)
			}()

			// // Вызов метода выполнения запроса select к СУБД ClickHouse
			// chs := make(chan string, 1)
			// chf := make(chan float64, 1)

			// wg.Add(1)
			// go func() error {

			// 	s := new(sl.CheckId)

			// 	s.Mu.Lock()
			// 	chs <- reqId.Id
			// 	s.Mu.Unlock()

			// 	s.Mu.Lock()
			// 	chf <- reqId.Measurement
			// 	s.Mu.Unlock()

			// 	err := s.SelectLog(wg, chs, chf)
			// 	if err != nil {
			// 		log.Printf("Error connect to dBase: %v", err)
			// 		return stream.SendMsg(err)
			// 	}
			// 	return nil
			// }()
			// go func() {
			// 	wg.Wait()
			// 	close(chs)
			// }()

			// Logic makes group of orders. Логика для объединения заказов в партии на основе адреса доставки
			destination := reqMap[reqId.GetId()].Id
			shipment, found := statusMap[destination]

			if found {
				req := reqMap[reqId.GetId()]
				shipment.IOList = append(shipment.IOList, &req)
				statusMap[destination] = shipment
			} else {
				comIO := pb.StatusIO{Id: "status - " + (reqMap[reqId.GetId()].Id), Status: "Processed!"}
				req := reqMap[reqId.GetId()]
				comIO.IOList = append(shipment.IOList, &req)
				statusMap[destination] = comIO
				log.Print(len(comIO.IOList), comIO.GetId())
			}

			if batchMarker == reqBatchSize {
				// Передаем клиенту поток заказов, объединенных в партии, group orderBatchSize
				for _, comb := range statusMap {
					// Group of orders. Передаем клиенту партию объединенных заказов
					log.Printf("Shipping : %v -> %v", comb.Id, len(comb.IOList))
					if err := stream.Send(&comb); err != nil { // Writes group of orders. Запись объединенных заказов в поток
						return err
					}
				}
				batchMarker = 0
				statusMap = make(map[string]pb.StatusIO)
			} else {
				batchMarker++
			}
		}
	}
}
