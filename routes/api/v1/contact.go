package v1

import (
	"net/http"
	"paper_back/exceptions"
	contact_service "paper_back/services/contact"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	contacts, err := contact_service.GetContacts(userId.(int64), pageSize, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, contacts)
}
