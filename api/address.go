package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sergeyignatov/simpleipam/common"
)

func apiGetIP(c *gin.Context) {
	subnet := c.PostForm("subnet")
	mac := c.PostForm("mac")
	fqdn := c.PostForm("fqdn")
	oldip := c.PostForm("ip")

	ip, err := context.Subnets.GetNewIp(subnet, mac, fqdn, oldip)
	if err != nil {
		Fail(c, err)
		return
	}
	c.JSON(200, common.NewApiResponse(ip))
}

func apiReleaseIP(c *gin.Context) {
	subnet := c.PostForm("subnet")
	mac := c.PostForm("mac")
	fqdn := c.PostForm("fqdn")
	ip := c.PostForm("ip")
	err := context.Subnets.ReleaseIP(subnet, mac, ip, fqdn)
	if err != nil {
		Fail(c, err)
		return
	}
	c.JSON(200, common.NewApiResponse("ok"))
}

func apiGetList(c *gin.Context) {
	c.JSON(200, context.Subnets.List())
}
