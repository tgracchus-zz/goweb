package main

import (
	"flag"
	"log"
	"os"
	"github.schibsted.io/spt-infrastructure/krakend/config"
	"github.schibsted.io/spt-infrastructure/krakend/config/viper"
	"github.schibsted.io/spt-infrastructure/krakend/logging"
	"github.schibsted.io/spt-infrastructure/krakend/logging/gologging"
	"github.schibsted.io/spt-infrastructure/krakend/proxy"
	gconfig "github.schibsted.io/spt-infrastructure/apigw-krakend/config"
	grouter "github.schibsted.io/spt-infrastructure/apigw-krakend/proxy/router/gin"
)

func main() {
	port := flag.Int("p", 0, "Port of the service")
	logLevel := flag.String("l", "ERROR", "Logging level")
	debug := flag.Bool("d", false, "Enable the debug")
	configFile := flag.String("c", "/etc/krakend/configuration.json", "Path to the configuration filename")
	flag.Parse()

	parser := viper.New()
	serviceConfig, err := parser.Parse(*configFile)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}
	serviceConfig.Debug = serviceConfig.Debug || *debug
	if *port != 0 {
		serviceConfig.Port = *port
	}

	logger, err := gologging.NewLogger(*logLevel, os.Stdout, "[KRAKEND]")
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}

	apiGWConfig := gconfig.NewApiGWConfig(*port, *debug, *configFile)
	gate := grouter.NewGateway(logger)
	gate.Run(apiGWConfig)
}

// customProxyFactory adds a logging middleware wrapping the internal factory
type customProxyFactory struct {
	logger  logging.Logger
	factory proxy.Factory
}

// New implements the Factory interface
func (cf customProxyFactory) New(cfg *config.EndpointConfig) (p proxy.Proxy, err error) {
	p, err = cf.factory.New(cfg)
	if err == nil {
		p = proxy.NewLoggingMiddleware(cf.logger, cfg.Endpoint)(p)
	}
	return
}
