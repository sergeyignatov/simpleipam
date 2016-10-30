package api

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyignatov/simpleipam/common"
	ctx "github.com/sergeyignatov/simpleipam/context"
	"net/http"
)

var context *ctx.Context

func Fail(c *gin.Context, err error) {
	c.Error(err)
	//c.JSON(500, gin.H{"Error": err.Error()})
	c.JSON(500, common.NewApiResponse(err))
}

func Router(c *ctx.Context) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	context = c
	router.Use(gin.ErrorLogger())
	//index := router.Group("/")
	root := router.Group("/api/1.0")
	{
		root.GET("/version", apiVersion)
		root.GET("/", apiGetList)
		root.POST("/getip", apiGetIP)
		root.POST("/release", apiReleaseIP)
	}

	return router
}
