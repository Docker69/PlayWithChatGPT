package router

import (
	mongodb "backend/db"
	mylogger "backend/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var apiKey = ""

// init Server
func init() {
}

// run server function
func RunServer(key string) {

	//set the OpenAI api key
	apiKey = key

	//get port from env
	port, exists := os.LookupEnv("LISTEN_PORT")

	if !exists {
		mylogger.Logger.Warn("LISTEN_PORT not defined in env, using default port 8080")
		port = "8080"
	}

	// create http server using Echo framework
	router := echo.New()

	// add recovery middleware
	router.Use(middleware.Recover())

	// add request logging middleware
	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			mylogger.Logger.WithFields(logrus.Fields{
				"uri":    values.URI,
				"status": values.Status,
				"method": values.Method,
			}).Info("Incoming request logged")

			return nil
		},
	}))

	// add a custom middleware to log the response
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// call the next middleware
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			// log the response
			mylogger.Logger.WithFields(logrus.Fields{
				"status": c.Response().Status,
				"method": c.Request().Method,
				"uri":    c.Request().RequestURI,
			}).Info("Response logged")

			return nil
		}
	})

	//TODO: fill actual id of the chat
	// add a custom middleware to set the response header
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// set the response header
			c.Response().Header().Set("X-Request-ID", uuid.New().String())

			// call the next middleware
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			return nil
		}
	})

	// Example API endpoint
	router.GET("/ping", handlePing)

	// Init Chat API endpoint
	router.POST("/api/v0/init/session", handleInitSession)

	// Init Chat API endpoint
	router.POST("/api/v0/init/chat", handleInitChat)

	// post chat completion to API endpoint
	router.POST("/api/v0/send-completion", handleChatCompletion)

	// get all chats list
	router.POST("/api/v0/getallchatslist", handleGetChatsList)

	//TODO: handle CORS properly
	//handle CORS
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			mylogger.Logger.Fatalf("listen: %s", err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 20 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	mylogger.Logger.Info("Shutting down server...")

	// The context is used to inform the server it has 20 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		if ctx.Err() != nil {
			// context already cancelled
			mylogger.Logger.Info("Server shutdown cancelled")
			return
		}
		mylogger.Logger.Fatalf("Server forced to shutdown: %s", err.Error())
	}

	if err := mongodb.Shutdown(ctx); err != nil {
		if ctx.Err() != nil {
			// context already cancelled
			mylogger.Logger.Info("MongoDB shutdown cancelled")
			return
		}
		mylogger.Logger.Fatalf("MongoDB forced to shutdown: %s", err.Error())
	}

	mylogger.Logger.Info("Server exiting")
}
