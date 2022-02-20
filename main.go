package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/strat0d/lvapi/lvget"
)

func main() {
	//GIN
	router := gin.Default()
	router.SetTrustedProxies(nil)

	ag_domains := router.Group("/api/v0/domains")
	{
		//get all domains on a host
		ag_domains.GET("/:host", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")

			res := make(chan lvget.DomainsResult)
			h := c.Param("host")
			go func() {
				res <- lvget.Domains(h)
				close(res)
			}()
			r := <-res
			if r.Error != nil {
				c.IndentedJSON(http.StatusOK, gin.H{"error": r.Error.Error})
			} else {
				c.IndentedJSON(http.StatusOK, r.Domains)
			}
		})
		//get domain :val :by(id, name, uuid) on :host
		ag_domains.GET("/:host/:by/:val", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")

			res := make(chan lvget.DomainResult)
			h, by, v := c.Param("host"), c.Param("by"), c.Param("val")
			go func() {
				res <- lvget.Domain(h, by, v)
				close(res)
			}()
			r := <-res
			if r.Error != nil {
				c.IndentedJSON(http.StatusOK, gin.H{"error": r.Error.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, r.Domain)
			}

		})

		//POST
		//run libvirt/virsh :action on guest :val found :by<id,name,uuid> on :host
		ag_domains.POST("/:host/:by/:val/:action", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")

			dom := make(chan lvget.LvDomainResult)
			h, by, v := c.Param("host"), c.Param("by"), c.Param("val")
			go func() {
				dom <- lvget.LvDomain(h, by, v, true)
				close(dom)
			}()
			d := <-dom
			if d.Error != nil {
				c.IndentedJSON(http.StatusOK, gin.H{"error": d.Error.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, d.Domain.Destroy())
			}

		})
	}

	ag_misc := router.Group("/api/v0/misc")
	{
		ag_misc.GET("/defaultxml", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			dom := lvget.DefaultXML()
			c.XML(http.StatusOK, dom)
		})
	}

	router.Run("0.0.0.0:8080")
}
