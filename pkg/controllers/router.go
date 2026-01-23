package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jroden2/holmes-go/pkg/controllers/public"
	"github.com/jroden2/holmes-go/pkg/middlewares"
	"github.com/rs/zerolog"
)

func Initialise(logger *zerolog.Logger) {
	router := gin.New()
	gin.SetMode(gin.DebugMode)

	HOST := os.Getenv("ServerHost")
	if HOST == "" {
		HOST = "localhost"
	}
	PORT := os.Getenv("ServerPort")
	if PORT == "" {
		PORT = "8080"
	}
	hostPort := fmt.Sprintf("%s:%s", HOST, PORT)
	logger.Info().Msg(hostPort)
	router.Use(gin.Recovery())
	router.Use(middlewares.ZerologMiddleware())

	router.Static("/css", "./templates/css")
	router.Static("/js", "./templates/js")
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "GET" {
			c.Redirect(302, "/")
			return
		}
		c.AbortWithStatus(404)
	})
	base := router.Group("")
	{
		public.Routes(base, logger)
	}

	srv := &http.Server{
		Addr:    hostPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error().Msgf("listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info().Msg("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Info().Msgf("Server Shutdown: %v", err)
	}

	logger.Info().Msg("Server exiting")
}
