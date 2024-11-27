package main

import (
	"context"
	"example.com/test/pkg/components/databases/dbconfig"
	"example.com/test/pkg/components/databases/repository"
	"example.com/test/pkg/config"
	"example.com/test/pkg/config/logger"
	"example.com/test/pkg/controllers"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func main() {
	config.LoadConfigurationMicroservice("./")
	logger, _ := logger.ApplyLoggerConfiguration(config.Configuration.Log.Level)
	defer closeLoggerHandler()(logger)

	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()
	RequestIDMiddleware()

	logger.Info("Starting server")

	// Commons services instances
	dbConfig := dbconfig.NewDatabaseConfig()
	pgConn, err := dbConfig.GetConnection()

	if err != nil {
		logger.Fatal("Error creating database connection")
		return
	}

	db := repository.New(pgConn)

	serverInstance := controllers.NewServer(db, pgConn, logger)

	stopCh := SetupSignalHandler()

	logger.Info("Creating routers")
	serverInstance.ListenAndServe(stopCh)
}

func closeLoggerHandler() func(logger *zap.Logger) {
	return func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func RequestIDMiddleware() {
	logger.Reset()

	requestID := uuid.New()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestID", requestID)

	logger.LoggerInstance = logger.LoggerInstance.With(zap.String("requestID", requestID.String()))
	logger.WithCtx(ctx, logger.LoggerInstance)
}

var onlyOneSignalHandler = make(chan struct{})

func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
