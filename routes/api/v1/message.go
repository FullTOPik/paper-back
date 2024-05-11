package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"paper_back/exceptions"
	message_model "paper_back/models/message"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addMessageBody struct {
	Text      string `json:"text"`
	ContactId int64  `json:"contactId"`
}

func AddMessage(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "Unauthorized",
		})
		return
	}

	jsonData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "error read body",
		})
		return
	}

	var bodyData addMessageBody

	if err = json.Unmarshal(jsonData, &bodyData); err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "error transform body json",
		})
		return
	}

	message, err := message_model.CreateMessage(userId.(int64), bodyData.ContactId, bodyData.Text)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, message)
}

func GetMessages(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "Unauthorized",
		})
		return
	}

	contactIdString, isPage := ctx.GetQuery("contactId")
	if !isPage {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "contactId is required",
		})
		return
	}

	contactId, err := strconv.Atoi(contactIdString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "contactId must be number",
		})
		return
	}

	messages, err := message_model.GetMessages(userId.(int64), int64(contactId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
