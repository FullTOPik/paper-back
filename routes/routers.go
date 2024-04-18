package routers

import (
	"paper_back/middlewares"
	v1 "paper_back/routes/api/v1"

	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/registration", v1.Registration)
	r.POST("/login", v1.Login)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(middlewares.CheckAuth())
	{

		apiv1.GET("/user/:id", v1.GetUser)
		apiv1.GET("contact", v1.GetContacts)
	}

	return r
}
