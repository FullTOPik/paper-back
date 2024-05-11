package routers

import (
	"paper_back/middlewares"
	v1 "paper_back/routes/api/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://localhost:5173"}
	r.Use(cors.New(config))

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("registration", v1.Registration)
	r.POST("login", v1.Login)

	apiv1 := r.Group("api/v1")
	apiv1.Use(middlewares.CheckAuth())
	{
		apiv1.GET("info", v1.GetInfo)
		apiv1.GET("user/:id", v1.GetUser)
		
		apiv1.GET("contact", v1.GetContacts)
		apiv1.POST("contact", v1.AddContact)
		apiv1.GET("code", v1.GetCode)

		apiv1.GET("message", v1.GetMessages)
		apiv1.POST("message", v1.AddMessage)
	}

	return r
}
