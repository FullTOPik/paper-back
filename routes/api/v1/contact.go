package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"paper_back/exceptions"
	contact_model "paper_back/models/contact"
	contact_service "paper_back/services/contact"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addContactBody struct {
	Code string `json:"code"`
}

func GetCode(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "Unauthorized",
		})
		return
	}

	code, err := contact_model.CreateUserCode(userId.(int64))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, code)
}

func AddContact(ctx *gin.Context) {
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

	var bodyData addContactBody

	if err = json.Unmarshal(jsonData, &bodyData); err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "error transform body json",
		})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "userId must be number",
		})
		return
	}

	contact, err := contact_service.CreateContact(userId.(int64), bodyData.Code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, contact)
}

func GetContacts(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "Unauthorized",
		})
		return
	}

	pageString, isPage := ctx.GetQuery("page")
	if !isPage {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "page is required",
		})
		return
	}

	pageSizeString, isPageSize := ctx.GetQuery("pageSize")
	if !isPageSize {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "pageSize is required",
		})
		return
	}

	page, err := strconv.Atoi(pageString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "page must be number",
		})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "pageSize must be number",
		})
		return
	}

	contacts, count, err := contact_service.GetContacts(userId.(int64), pageSize, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, struct {
		Data  []contact_model.Contact `json:"data"`
		Count int64                   `json:"count"`
	}{
		Data:  contacts,
		Count: count,
	})
}
