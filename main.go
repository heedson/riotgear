package main

import (
	"context"
	"fmt"
	"mime"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kelseyhightower/envconfig"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/heedson/riotgear/api"
	"github.com/heedson/riotgear/proto"
	_ "github.com/heedson/riotgear/statik"
)

func mustParseURL(serverName string) (serverURL *url.URL) {
	serverURL, err := url.Parse(fmt.Sprintf("https://%s.api.riotgames.com", serverName))
	if err != nil {
		panic(err)
	}

	return serverURL
}

// serveOpenAPI serves an OpenAPI UI on /openapi-ui/
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
func serveOpenAPI(mux *http.ServeMux) error {
	mime.AddExtensionType(".svg", "image/svg+xml")

	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	// Expose files in static on <host>/openapi-ui
	fileServer := http.FileServer(statikFS)
	prefix := "/openapi-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	return nil
}

type config struct {
	RiotAPIKey  string `required:"true" envconfig:"RIOT_API_KEY" desc:"The Riot API key to use for access to the Riot API."`
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

	var envOpts config
	if err := envconfig.Process("", &envOpts); err != nil {
		envconfig.Usage("", &envOpts)
		logger.WithError(err).Fatal()
	}

	var regionsToServerURL = map[string]*url.URL{
		"br":   mustParseURL("br1"),
		"eune": mustParseURL("eun1"),
		"euw":  mustParseURL("euw1"),
		"jp":   mustParseURL("jp1"),
		"kr":   mustParseURL("kr"),
		"lan":  mustParseURL("la1"),
		"las":  mustParseURL("la2"),
		"na":   mustParseURL("na1"),
		"oce":  mustParseURL("oc1"),
		"tr":   mustParseURL("tr1"),
		"ru":   mustParseURL("ru"),
		"pbe":  mustParseURL("pbe1"),
	}

	s := grpc.NewServer()

	srv := api.NewServer(
		logger,
		&http.Client{
			Timeout: time.Second * 10,
		},
		regionsToServerURL,
		envOpts.RiotAPIKey,
	)

	proto.RegisterRiotgearServer(s, srv)

	go func() {
		lis, err := net.Listen("tcp", envOpts.GRPCAddr)
		if err != nil {
			logger.WithError(err).Fatal("Failed to start grpc listener")
		}

		if err = s.Serve(lis); err != nil {
			logger.WithError(err).Fatal("Failed to serve gRPC server")
		}
	}()

	con, err := grpc.Dial(envOpts.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		logger.WithError(err).Fatal("Failed to dial gRPC server")
	}

	mux := http.NewServeMux()
	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption("*", &runtime.JSONPb{
			EmitDefaults: true,
		}),
	)
	if err = proto.RegisterRiotgearHandler(context.Background(), gwMux, con); err != nil {
		logger.WithError(err).Fatal("Failed to register riotgear in gRPC-gateway")
	}

	mux.Handle("/", gwMux)

	if err = serveOpenAPI(mux); err != nil {
		logger.WithError(err).Fatal("Failed to serve OpenAPI UI")
	}

	logger.Infof("Serving gRPC-Gateway on http://%s", envOpts.GatewayAddr)
	logger.Infof("Serving OpenAPI Documentation on http://%s/openapi-ui/", envOpts.GatewayAddr)

	gwServer := http.Server{
		Addr:    envOpts.GatewayAddr,
		Handler: mux,
	}

	if err = gwServer.ListenAndServe(); err != nil {
		logger.WithError(err).Fatal("Failed to serve gRPC-gateway")
	}
}
