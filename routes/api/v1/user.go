package v1

import (
	"net/http"
	"paper_back/exceptions"
	user_service "paper_back/services/user"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	LastVisit time.Time `json:"last_visit"`
	Password  string    `json:"password"`
}

func GetInfo(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "Unauthorized",
		})
		return
	}

	user, err := user_service.GetUser(userId.(int64))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func GetUser(ctx *gin.Context) {
	userIdString := ctx.Param("id")

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: "Id must be number",
		})
		return
	}

	user, err := user_service.GetUser(int64(userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func Registration(ctx *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	token, err := user_service.Registration(body.Username, body.Password, "USER")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.Header("token", token)
	ctx.SetCookie("token", token, int(time.Hour)/1000000, "*", "*", false, true)
	ctx.Status(http.StatusOK)
}

func Login(ctx *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	token, err := user_service.Login(body.Username, body.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.ErrorWithStatus{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	ctx.SetCookie("token", token, int(time.Hour)/1000000, "*", "*", false, false)
	ctx.Status(http.StatusOK)
}
