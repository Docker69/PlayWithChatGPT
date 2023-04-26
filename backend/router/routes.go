package router

import (
	"backend/ai"
	"backend/db/mongodb"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

func handlePing(c echo.Context) error {

	// log a message using logrus logger
	utils.Logger.Info("Received request")

	return c.JSON(http.StatusOK, map[string]string{
		"message": "pong",
	})
}

func handleInitSession(c echo.Context) error {

	// log a message using logrus logger
	utils.Logger.Info("Received session init request")

	// init request to empty JSON
	reqBody := models.Human{}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	//get nickname from request body
	if reqBody.NickName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Empty nickname",
		})
	}

	//find the user in the database
	human, err := mongodb.HumansCollection.GetByNickname(reqBody.NickName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if human.Id == "" {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Not Found",
		})
	}
	// return reqBody as json
	return c.JSON(http.StatusOK, human)

}

func handleInitChat(c echo.Context) error {

	// log a message using logrus logger
	utils.Logger.Info("Received chat init request")

	// get request body
	reqBody := models.ChatCompletionRequestBody{}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// add the message to a list of messages
	reqBody.Messages = append(reqBody.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: reqBody.Role,
	})

	var err error = nil
	//insert into DB the chat
	reqBody.Id, err = mongodb.ChatsCollection.Insert(&reqBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	//find the human in the database
	human, err := mongodb.HumansCollection.GetById(reqBody.HumanId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	//add chat id and role to human
	chatRecord := models.ChatRecord{
		Id:   reqBody.Id,
		Role: reqBody.Role,
	}
	human.ChatIds = append(human.ChatIds, chatRecord)

	//update the human in the database
	err = mongodb.HumansCollection.UpdateChats(&human)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// return reqBody as json
	return c.JSON(http.StatusOK, reqBody)

}

// get chat by is
func handleGetChatById(c echo.Context) error {

	// log a message using logrus logger
	utils.Logger.Info("Received Get Chat request")

	//get chat id from request
	chatId := c.Param("id")

	//find the chat in the database
	chat, err := mongodb.ChatsCollection.GetById(chatId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if chat.Id == "" {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Not Found",
		})
	}
	// return reqBody as json
	return c.JSON(http.StatusOK, chat)

}

func handleChatCompletion(c echo.Context) error {

	// log a message using logrus logger
	utils.Logger.Info("Received chat completion request")

	// get request body
	reqBody := models.ChatCompletionRequestBody{}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	//Call the chat completion function and get the response, handle error
	resp, err := ai.ChatCompletion(reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	//spread request body and replace Message with response
	reqBody.Messages = resp

	//update the chat in the database
	err = mongodb.ChatsCollection.Update(&reqBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error updating DB": err.Error(),
		})
	}

	// return reqBody as json
	return c.JSON(http.StatusOK, reqBody)

}

func handleGetChatsList(c echo.Context) error {

	// log a message using logrus logger
	utils.Logger.Info("Received Get All Chats request")

	//Call the chat completion function and get the response, handle error
	resBody, err := mongodb.ChatsCollection.GetAll()

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// return reqBody as json
	return c.JSON(http.StatusOK, resBody)
}
