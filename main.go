package main

import (
	"github.com/appleboy/gin-jwt"
	"github.com/claywong/gin_api_demo/apps"
	"github.com/claywong/gin_api_demo/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/lunny/log"
	"net/http"
	"time"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

func helloHandler(ctx *gin.Context) {
	claim := jwt.ExtractClaims(ctx)
	user, _ := ctx.Get(identityKey)
	ctx.JSON(http.StatusOK, gin.H{
		"userId":   claim["id"],
		"UserName": user.(*User).UserName,
		"text":     "hello world.",
	})

}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.StaticFS("/static", http.Dir("static"))
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			// use pusher.Push() to do server push
			if err := pusher.Push("/static/layui/css/layui.css", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
			if err := pusher.Push("/static/css/base.css", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		c.HTML(http.StatusOK, "default.html", gin.H{
			"title": "API-DEMO-GIN",
		})
	})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(context *gin.Context) interface{} {
			claim := jwt.ExtractClaims(context)
			return &User{
				UserName: claim["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &User{
					UserName:  userID,
					LastName:  "Bo-Yi",
					FirstName: "Wu",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT ERROR" + err.Error())
	}

	router.POST("/login", authMiddleware.LoginHandler)
	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claim: %#v\n", claims)
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "page not found",
		})

	})

	auth := router.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", helloHandler)
	}

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

	router.RunTLS(":8000", "./testdata/server.pem","./testdata/server.key")
}
