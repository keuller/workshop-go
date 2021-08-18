package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/keuller/exchange/internal/infra"
	"github.com/ugorji/go/codec"
)

type App struct {
	server    *http.Server
	rpcServer rpc.ServerCodec
}

func New() App {
	host := infra.GetConfig("host")
	port := infra.GetConfig("port")

	routes := configureRoutes()

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      routes,
		IdleTimeout:  15 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return App{srv, nil}
}

func (a App) StartRpc() {
	var handle codec.MsgpackHandle
	listener, err1 := net.Listen("tcp", ":8005")
	if err1 != nil {
		log.Fatalf("[ERROR] cannot init RPC server, reason: %s \n", err1.Error())
	}

	log.Printf("[INFO] RPC Server is ready on tcp://0.0.0.0:8005 \n")
	for {
		conn, err2 := listener.Accept()
		if err2 != nil {
			log.Fatalf("[ERROR] cannot read RPC message, reason: %s \n", err2.Error())
		}

		a.rpcServer = codec.MsgpackSpecRpc.ServerCodec(conn, &handle)
		rpc.ServeCodec(a.rpcServer)
	}
}

func (a App) Start() {
	log.Printf("[INFO] Server is ready on http://%s\n", a.server.Addr)

	err := a.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("[ERROR] Cannot initiate the server, reason: %s \n", err.Error())
	}
}

func (a App) Server() *http.Server {
	return a.server
}

func (a App) Stop() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ctx.Done()

	log.Println("[INFO] Server is shutting down...")
	// stop receiving any request.
	if err := a.server.Shutdown(context.Background()); err != nil {
		log.Fatalf("[FAIL] Fail on stop the server, reason: %s\n", err.Error())
	}
	log.Println("[INFO] Server has been stopped.")
	stop()

	a.rpcServer.Close()
}
