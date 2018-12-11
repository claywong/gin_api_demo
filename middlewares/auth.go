package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/go/src/pkg/bytes"
	"github.com/golang/go/src/pkg/crypto/md5"
	"github.com/golang/go/src/pkg/fmt"
	"github.com/golang/go/src/pkg/net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		appkey := c.Request.FormValue("appkey")
		sign := c.Request.FormValue("sign")
		ts := c.Request.FormValue("ts")
		method := c.Request.Method

		if appkey != "1001" {
			noAUth(c, "key not found")
			return
		}


		if Sign(appkey, ts, method, sign) == false {
			noAUth(c, "Unauthorized")
			return
		}
		c.Next()
		return

	}
}

func noAUth(c *gin.Context, msg string) {

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": msg,
	})
	c.Abort()
}

func Sign(appkey, ts, method, sign string) bool {
	secret := "secret"

	b := bytes.Buffer{}
	b.WriteString("appkey=")
	b.WriteString(appkey)
	b.WriteString("&app_secret=")
	b.WriteString(secret)
	b.WriteString("&method=")
	b.WriteString(method)
	if Md5([]byte(b.String())) == sign {
		return true
	}
	return false

}

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
