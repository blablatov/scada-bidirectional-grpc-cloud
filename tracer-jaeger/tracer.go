package tracer

import (
	"log"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

func InitTracer() {

	// Тип экспортера Jaeger, путь для сбора данных, имя сервиса и конечная точка агента
	// Exporter type, path and endpoint
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	agentEndpointURI := "localhost:6831"
	collectorEndpointURI := "http://localhost:14268/api/traces"
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: collectorEndpointURI,
		AgentEndpoint:     agentEndpointURI,
		ServiceName:       "grpc-cloud",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Регистрация экспортера в трассировщике OpenCensus.
	// Registering exporter to tracer
	trace.RegisterExporter(exporter)
}
