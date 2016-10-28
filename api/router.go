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
	router := gin.Default()
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
	/*{
		root.GET("/containers", apiContainersList)
		root.POST("/containers", apiContainerCreate)
		root.GET("/containers/:name", apiContainersShow)
		root.POST("/containers/:name", apiContainersEdit)

		root.DELETE("/containers/:name", apiContainersDelete)
		root.POST("/containers/:name/start", apiContainersStart)
		root.POST("/containers/:name/stop", apiContainersStop)
	}
	{
		root.GET("/servers", apiServersList)
		root.POST("/servers/:name/maint", apiServersMaintenance)
	}
	{
		root.GET("/profiles", apiProfilesList)
		root.GET("/profiles/:name", apiProfilesShow)
		root.POST("/profiles/:name", apiProfilesEdit)
		root.DELETE("/profiles/:name", apiProfilesDelete)
	}
	metadata := router.Group("/metadata/2009-04-04/meta-data")
	{
		metadata.GET("/instance-id", func(c *gin.Context) {
			fmt.Println(c.Request)
			c.JSON(200, "fairsnail2")
		})
		metadata.GET("/", func(c *gin.Context) {
			fmt.Println(c.Request)
			c.JSON(200, "ok")
		})
	}
	*/
	return router
}
