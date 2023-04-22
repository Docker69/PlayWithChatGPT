package main

import (
	"backend/chat"
	mongodb "backend/db"
	mylogger "backend/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	completionmodels "backend/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// run server function
func runServer(apiKey string) {
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

	/* 		// use mylogger middleware
	   		router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	   			Format: "${status} ${method} ${uri} ${latency_human}\n",
	   			Output: mylogger.Logger.Out,
	   		}))

	   		// set response logging to JSON
	   		router.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	   			type ResponseLog struct {
	   				Status     int    `json:"status"`
	   				StatusText string `json:"statustext,omitempty"`
	   				//Data       interface{} `json:"data,omitempty"`
	   			}
	   			id := c.Response().Header().Get(echo.HeaderXRequestID)
	   			logData := ResponseLog{
	   				Status:     c.Response().Status,
	   				StatusText: http.StatusText(c.Response().Status),
	   				//Data:       c.Get("response"),
	   			}
	   			//msg := fmt.Sprintf(`{"reqid":"%v", "response":%+v}`, id, logData)
	   			mylogger.Logger.Infof(`{"reqid":"%v", "response":%+v}`, id, logData)
	   		}))
	*/
	// Example API endpoint
	router.GET("/ping", func(c echo.Context) error {

		// log a message using logrus logger
		mylogger.Logger.Info("Received request")

		return c.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	// Init Chat API endpoint
	router.POST("/api/init", func(c echo.Context) error {

		// log a message using logrus logger
		mylogger.Logger.Info("Received chat init request")

		// get request body
		reqBody := completionmodels.ChatCompletionRequestBody{}

		if err := c.Bind(&reqBody); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid request body",
			})
		}

		// generate a random uuid
		reqBody.Id = uuid.New().String()

		// return reqBody as json
		return c.JSON(http.StatusOK, reqBody)

	})

	// Init Chat API endpoint
	router.POST("/api/send-completion", func(c echo.Context) error {

		// log a message using logrus logger
		mylogger.Logger.Info("Received chat completion request")

		// get request body
		reqBody := completionmodels.ChatCompletionRequestBody{}

		if err := c.Bind(&reqBody); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid request body",
			})
		}

		//Call the chat completion function and get the response, handle error
		resp, err := chat.ChatCompletion(apiKey, reqBody)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		//spread request body and replace Message with response
		reqBody.Messages = resp

		// return reqBody as json
		return c.JSON(http.StatusOK, reqBody)

	})

	// Init Chat API endpoint
	router.POST("/api/getallchats", func(c echo.Context) error {

		// log a message using logrus logger
		mylogger.Logger.Info("Received Get All Chats request")

		//Call the chat completion function and get the response, handle error
		resBody, err := mongodb.GetAllChats()

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		// return reqBody as json
		return c.JSON(http.StatusOK, resBody)
	})

	//handle CORS
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	//get port from env
	port, exists := os.LookupEnv("LISTEN_PORT")

	if !exists {
		mylogger.Logger.Warn("LISTEN_PORT not defined in env, using default port 8080")
		port = "8080"
	}

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

	mylogger.Logger.Info("Server exiting")
}
