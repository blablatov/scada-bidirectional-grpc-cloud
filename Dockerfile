FROM golang:1.20

RUN git clone https://github.com/blablatov/scada-bidirectional-grpc-cloud.git
WORKDIR scada-bidirectional-grpc-cloud/grpc-service-cloud

RUN go mod download

COPY *.go ./
COPY *.conf ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /grpc-service-cloud
EXPOSE 50051

#WORKDIR scada-bidirectional-grpc-cloud/grpc-client-io
#RUN CGO_ENABLED=0 GOOS=linux go build -o /grpc-client-io

CMD ["/grpc-service-cloud"]
