package router

import (
	//	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	//	"backend/db/mongodb"
	//	"backend/models"
)

func TestHandlePing(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlePing(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var respMap map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &respMap)
	assert.NoError(t, err)
	assert.Equal(t, "pong", respMap["message"])
}

/*
func TestHandleInitSessionBadRequest(t *testing.T) {
	reqBody := []byte(`{"invalid": "request body"}`)
	req := httptest.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()
	c := router.CreateContext(req, rec)

	err := router.HandleInitSession(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var respMap map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &respMap)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request body", respMap["error"])
}

func TestHandleInitSessionEmptyNickname(t *testing.T) {
	reqBody := models.Human{NickName: ""}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(reqBytes))
	rec := httptest.NewRecorder()
	c := router.CreateContext(req, rec)

	err := router.HandleInitSession(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var respMap map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &respMap)
	assert.NoError(t, err)
	assert.Equal(t, "Empty nickname", respMap["error"])
}

func TestHandleInitSessionNotFound(t *testing.T) {
	nickname := "nonexistent"
	reqBody := models.Human{NickName: nickname}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(reqBytes))
	rec := httptest.NewRecorder()
	c := router.CreateContext(req, rec)
	mongodb.HumansCollection.Insert(&models.Human{NickName: "existing"})

	err := router.HandleInitSession(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	var respMap map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &respMap)
	assert.NoError(t, err)
	assert.Equal(t, "Not Found", respMap["error"])
}

func TestHandleInitSessionSuccess(t *testing.T) {
	nickname := "human123"
	reqBody := models.Human{NickName: nickname}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/session", bytes.NewBuffer(reqBytes))
	rec := httptest.NewRecorder()
	c := router.CreateContext(req, rec)
	human := models.Human{NickName: nickname}
	mongodb.HumansCollection.Insert(&human)

	err := router.HandleInitSession(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var respBody models.Human
	err = json.Unmarshal(rec.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.Equal(t, human.Id, respBody.Id)
	assert.Equal(t, human.NickName, respBody.NickName)
}

func TestHandleInitChatBadRequest(t *testing.T) {
	reqBody := []byte(`{"invalid": "request body"}`)
	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()
	c := router.CreateContext(req, rec)

	err := router.HandleInitChat(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var respMap map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &respMap)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request body", respMap["error"])
}

func TestHandleInitChatDBError(t *testing.T) {
	humanID := "nonexistent"
	role := "customer"
	reqBody := models.ChatCompletionRequestBody{HumanId: humanID, Role: role}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(reqBytes))
	rec := httptest.NewRecorder()
	c := router.CreateContext(req, rec)
	mongodb.HumansCollection.Insert(&models.Human{NickName: "existing"})

	err := router.HandleInitChat(c)
	assert.Error(t, err) // assert error returned
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var respMap map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &respMap)
	assert.NoError(t, err)
	assert.Contains(t, respMap["error"], "not found")
}

func TestHandleInitChatSuccess(t *testing.T) {
	humanID := "human123"
	role := "customer"
	reqBody := models.ChatCompletionRequestBody{HumanId: humanID, Role: role}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(reqBytes))
	rec := httptest.NewRecorder()
	c := router.CreateContext(req, rec)
	human := models.Human{Id: humanID, NickName: "human"}
	mongodb.HumansCollection.Insert(&human)

	err := router.HandleInitChat(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var respBody models.ChatCompletionRequestBody
	err = json.Unmarshal(rec.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, respBody.Id)
	assert.Equal(t, humanID, respBody.HumanId)

	// check that chat record was added to human's chatIds and saved in DB
	humanFromDB, _ := mongodb.HumansCollection.GetById(humanID)
	assert.Len(t, humanFromDB.ChatIds, 1)
	assert.Equal(t, respBody.Id, humanFromDB.ChatIds[0].Id)
}
*/
