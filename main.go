package main

import (
	"github.com/claywong/gin_api_demo/apps"
	"github.com/claywong/gin_api_demo/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.StaticFS("static", http.Dir("static"))
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "default.html", gin.H{
			"title": "API-DEMO-GIN",
		})
	})

	v0 := router.Group("/v0")
	v0.Use(middlewares.Auth())
	{

		//新增
		//curl -X POST http://127.0.0.1:8000/v0/member -d "login_name=hello&password=1234"
		v0.POST("/member", apps.MemberAdd)
		// curl -X GET http://127.0.0.1:8000/v0/member
		v0.GET("/member", apps.MemberList)
		// curl -X GET http://127.0.0.1:8000/v0/member/1
		v0.GET("/member/:id", apps.MemberGet)
		//curl -X PUT http://127.0.0.1:8000/v0/member/1 -d "login_name=hello&password=1234"
		v0.PUT("/member/:id", apps.MemberEdit)
		// curl -X DELETE http://127.0.0.1:8000/v0/member/2
		v0.DELETE("/member/:id", apps.MemberDelete)
	}

	router.Run(":8000")
}
