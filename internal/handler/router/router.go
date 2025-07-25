package router

import (
	"net/http"

	_ "people/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router(server *Server) *gin.Engine {
	router := gin.Default()
	handler := New(server.usecase, server.log)
	api := router.Group("/api/v1")
	{
		api.GET("/swagger", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/people/swagger/index.html")
		})
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		api.GET("/users/:id", handler.GetUserInfoBySecondName)
		api.GET("/users", handler.GetAllUsersInfo)
		api.GET("/users/:id/emails", handler.GetUserEmails)
		api.GET("/users/:id/friends", handler.GetUserFriends)
		api.POST("/users", handler.CreateUser)
		api.POST("/users/:id/emails", handler.AddUserEmails)
		api.POST("/users/:id/friends", handler.AddUserFriends)
		api.PUT("/users/:id", handler.UpdateUser)
		api.DELETE("/users/:id", handler.DeleteUser)
		api.DELETE("/users/emails", handler.DeleteEmails)
		api.DELETE("/users/:id/friends", handler.DeleteUserFriends)
	}
	return router
}
