package main

import (
	"fmt"
	"log"

	"suggest-runtime/internal/config"
	serverImpl "suggest-runtime/internal/server"
	"suggest-runtime/internal/util/pprof"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.LoadConfig("config-local.json")
	if err != nil {
		log.Default().Fatal(err)
		return
	}

	if cfg.Pprof.Enable {
		pprof.HandlePprof()
	}

	server := serverImpl.NewServer(cfg)

	e := echo.New()
	e.Use(middleware.CORS())

	serverImpl.RegisterHandlers(e, server)

	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	err = e.Start(address)
	log.Default().Fatal(err)
}
