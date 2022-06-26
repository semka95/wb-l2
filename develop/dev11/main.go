package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"calendar/event/api"
	"calendar/event/repository/bolt"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	logger.Info("reading config")
	config, err := NewConfig()
	if err != nil {
		logger.Error("can't decode config", zap.Error(err))
		return
	}

	logger.Info("connecting to database")
	db, err := bolt.NewBoltDB(config.DBPath)
	if err != nil {
		panic(err)
	}

	store := bolt.NewBoltEventRepository(db)
	api := api.NewAPI(store, logger)
	router := api.NewRouter()

	srv := &http.Server{
		Addr:        config.HTTPServerAddress,
		Handler:     router,
		ReadTimeout: time.Duration(config.ReadTimeout) * time.Second,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
	}

	logger.Info("running http server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("can't start server", zap.Error(err), zap.String("server address", config.HTTPServerAddress))
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info("received interrupt signal, closing server")
	timeout, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	if err := srv.Shutdown(timeout); err != nil {
		logger.Error("can't shutdown http server", zap.Error(err))
	}
}
