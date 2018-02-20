package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/heedson/riotgear/api"
	"github.com/heedson/riotgear/proto"
)

type config struct {
	GRPCAddr    string `default:"localhost:8081" envconfig:"GRPC_ADDR" desc:"Address to serve the gRPC Server on."`
	GatewayAddr string `default:"0.0.0.0:8080" split_words:"true" desc:"Address to serve the gRPC-Gateway on."`
}

func main() {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: time.StampMilli,
		FullTimestamp:   true,
	}

	envOpts := config{}
	err := envconfig.Process("", &envOpts)
	if err != nil {
		envconfig.Usage("", &envOpts)
		logger.WithError(err).Fatal()
	}

	s := grpc.NewServer()

	srv := api.NewServer(logger)

	proto.RegisterEchoTestServer(s, srv)

	go func() {
		lis, err := net.Listen("tcp", envOpts.GRPCAddr)
		if err != nil {
			logger.WithError(err).Fatal("Failed to start grpc listener")
		}

		err = s.Serve(lis)
		if err != nil {
			logger.WithError(err).Fatal("Failed to serve gRPC server")
		}
	}()

	cc, err := grpc.Dial(envOpts.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		logger.WithError(err).Fatal("Failed to dial gRPC server")
	}

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption("*", &runtime.JSONPb{
			EmitDefaults: true,
		}),
	)
	err = proto.RegisterEchoTestHandler(context.Background(), mux, cc)
	if err != nil {
		logger.WithError(err).Fatal("Failed to register echo test in gRPC-gateway")
	}

	logger.Infoln("Serving on", envOpts.GatewayAddr)
	err = http.ListenAndServe(envOpts.GatewayAddr, mux)
	if err != nil {
		logger.WithError(err).Fatal("Failed to serve gRPC-gateway")
	}
}
