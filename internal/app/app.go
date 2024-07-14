// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/maxyong7/chat-messaging-app/config"
	v1 "github.com/maxyong7/chat-messaging-app/internal/controller/http/v1"
	"github.com/maxyong7/chat-messaging-app/internal/usecase"
	"github.com/maxyong7/chat-messaging-app/internal/usecase/repo"
	"github.com/maxyong7/chat-messaging-app/internal/usecase/webapi"
	"github.com/maxyong7/chat-messaging-app/pkg/httpserver"
	"github.com/maxyong7/chat-messaging-app/pkg/logger"
	"github.com/maxyong7/chat-messaging-app/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.PostgresURL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	userInfoRepo := repo.NewUserInfo(pg)
	// Use case
	translationUseCase := usecase.New(
		repo.New(pg),
		webapi.New(),
	)
	verificationUseCase := usecase.NewAuth(
		userInfoRepo,
		webapi.New(),
	)
	conversationUseCase := usecase.NewConversation(
		repo.NewConversation(pg),
		userInfoRepo,
	)
	inboxUseCase := usecase.NewInbox(
		repo.NewConversation(pg),
	)
	contactUseCase := usecase.NewContacts(
		repo.NewContacts(pg),
		userInfoRepo,
	)

	// // RabbitMQ RPC Server
	// rmqRouter := amqprpc.NewRouter(translationUseCase)

	// rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	// if err != nil {
	// 	l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	// }

	// HTTP Server
	handler := gin.New()
	routerUseCase := v1.RouterUseCases{
		Translation:  translationUseCase,
		Verification: verificationUseCase,
		Conversation: conversationUseCase,
		Inbox:        inboxUseCase,
		Contact:      contactUseCase,
	}
	v1.NewRouter(handler, l, routerUseCase)

	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		// case err = <-rmqServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	// err = rmqServer.Shutdown()
	// if err != nil {
	// 	l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	// }
}
