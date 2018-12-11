package apps

import (
	"github.com/claywong/gin_api_demo/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/go/src/pkg/log"
	"github.com/golang/go/src/pkg/net/http"
	"github.com/golang/go/src/pkg/strconv"
)

func MemberList(ctx *gin.Context) {

	filters := make([]interface{}, 0)
	filters = append(filters, "id", "<>", "0")

	page, _ := strconv.Atoi(ctx.Request.FormValue("page"))
	pageSize, _ := strconv.Atoi(ctx.Request.FormValue("page_size"))

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 10
	}

	list, n, err := models.ListMember(page, pageSize, filters...)

	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"status":  http.StatusExpectationFailed,
			"message": err.Error(),
			"data":    "123",
		})
		log.Fatal(err)
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"message":   "SUCCESS",
			"data":      list,
			"count":     n,
			"page_size": pageSize,
			"current":   page,
		})
	}

}

func MemberGet(ctx *gin.Context) {

}

func MemberEdit(ctx *gin.Context) {

}

func MemberAdd(ctx *gin.Context) {

}

func MemberDelete(ctx *gin.Context) {

}
