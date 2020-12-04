package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/byuoitav/crestron/cpu3/matrix"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/pflag"
)

func main() {
	var (
		port int
	)

	pflag.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.Parse()

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err)
		os.Exit(1)
	}

	switchers := &sync.Map{}

	handlers := Handlers{
		CreateVideoSwitcher: func(addr string) *matrix.Matrix {
			if vs, ok := switchers.Load(addr); ok {
				return vs.(*matrix.Matrix)
			}

			vs := matrix.New(addr)

			// these numbers have been tested for the 16x16
			vs.OutputSlotStart = 33
			vs.SetRouteOutputStart = 101

			switchers.Store(addr, &vs)
			return &vs
		},
	}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	api := e.Group("/api/v1")
	handlers.RegisterRoutes(api)

	log.Printf("Server started on %v", lis.Addr())
	if err := e.Server.Serve(lis); err != nil {
		log.Printf("unable to serve: %s", err)
	}
}
