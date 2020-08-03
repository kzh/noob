package tracing

import (
	"io"
	"log"

	opentracing "github.com/opentracing/opentracing-go"
	config "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func InitJaeger() (io.Closer, error) {
	log.Println("Connecting to Jaeger...")

	cfg, err := config.FromEnv()
	if err != nil {
		log.Printf("Could not parse Jaeger env vars: %s", err.Error())
		return nil, err
	}

	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaegerlog.StdLogger),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	log.Println("Connected to Jaeger.")
	return closer, nil
}
