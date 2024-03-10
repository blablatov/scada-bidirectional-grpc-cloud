### Тестирование клиентского кода без подключения к серверу. 
### Testing code without connect to server              

### Go 1.16+
```shell script
go install github.com/golang/mock/mockgen@v1.6.0
```

### Использование Gomock, генерация макетов интерфейсов клиентского gRPC-приложения. Use Gomock      
Для генерации макета интерфейса CloudExchangeClient, перейти в `github.com/blablatov/scada-bidirectional-grpc-cloud/grpc-mockups` и выполнить       
(Runs generation code of mock up for interface CloudExchangeClient): 
       
```shell script
mockgen github.com/blablatov/scada-bidirectional-grpc-cloud/grpc-cloud-proto CloudExchangeClient > mock_grpc_client.go
```

### Тестирование. Run test    

```shell script
go build . && go run .
```


