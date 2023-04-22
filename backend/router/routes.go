package router

import (
	"backend/chat"
	mongodb "backend/db"
	completionmodels "backend/models"
	mylogger "backend/utils"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func handlePing(c echo.Context) error {

	// log a message using logrus logger
	mylogger.Logger.Info("Received request")

	return c.JSON(http.StatusOK, map[string]string{
		"message": "pong",
	})
}

func handleInitChat(c echo.Context) error {

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

}

func handleChatCompletion(c echo.Context) error {

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

}

func handleGetChatsList(c echo.Context) error {

	// log a message using logrus logger
	mylogger.Logger.Info("Received Get All Chats request")

	//Call the chat completion function and get the response, handle error
	resBody, err := mongodb.ChatsCollection.GetAll(context.TODO())

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// return reqBody as json
	return c.JSON(http.StatusOK, resBody)
}
