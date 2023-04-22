package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/cmd/urlshortener/env"
	"url-shortener/cmd/urlshortener/internal/urlshortener"
	"url-shortener/pkg/encoder"
	"url-shortener/pkg/repository/firestore/counter"
	"url-shortener/pkg/repository/firestore/urls"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const shardsNumber = 100

func main() {
	logrus.Info("loading application config...")
	config, err := env.LoadAppConfig()
	if err != nil {
		logrus.Fatalf("failed to create firestore client: ", err)
	}

	logrus.Info("establishing firestore connection...")
	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		logrus.Fatalf("failed to create firestore client: ", err)
	}

	urlsRepository := urls.NewRepository(firestoreClient)
	counterRepository := counter.NewRepository(firestoreClient, shardsNumber)
	controller := urlshortener.NewController(urlsRepository, counterRepository, encoder.New())
	presenter := urlshortener.NewPresenter(controller)

	logrus.Info("initializing shards...")
	if err := counterRepository.InitCounter(ctx); err != nil {
		logrus.Fatal("failed to initialize counter: ", err)
	}

	handler := gin.Default()
	handler.POST("/", presenter.CreateShortURL)
	handler.GET("/:short_url", presenter.RedirectToLongURL)

	logrus.Info("http server is starting...")
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: handler,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Fatal("failed to listen and serve: ", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigChan
	signal.Stop(sigChan)
	logrus.Info("http server is stopping...")

	shutdownCtx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFunc()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logrus.Fatal("failed to shutdown server", err)
	}
}
