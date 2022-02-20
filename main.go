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
			res := make(chan lvget.DomainsResult)
			c.Header("Access-Control-Allow-Origin", "*")
			cc := c.Copy()
			go func(c *gin.Context) {
				res <- lvget.GetDomains(cc)
				close(res)
			}(cc)
			r := <-res
			if r.Err != nil {
				c.IndentedJSON(http.StatusOK, gin.H{"error": r.Err.Error})
			} else {
				c.IndentedJSON(http.StatusOK, r.Domains)
			}
		})
		//get domain :val by :<method>(id, name, uuid) on :host
		ag_domains.GET("/:host/:method/:val", func(c *gin.Context) {
			res := make(chan lvget.DomainResult)
			c.Header("Access-Control-Allow-Origin", "*")
			cc := c.Copy()
			go func(c *gin.Context) {
				res <- lvget.GetDomain(c)
				close(res)
			}(cc)
			r := <-res
			if r.Err != nil {
				c.IndentedJSON(http.StatusOK, gin.H{"error": r.Err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, r.Domain)
			}

		})

		//POST
		//run libvirt/virsh :action on guest :val found by :method<id,name,uuid> on :host
		ag_domains.POST("/:host/:method/:val/:action", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			//
		})
	}

	ag_misc := router.Group("/api/v0/misc")
	{
		ag_misc.GET("/defaultxml", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			dom := lvget.GetDefaultXML()
			c.XML(http.StatusOK, dom)
		})
	}

	router.Run("0.0.0.0:8080")
}

/* func getVersion(l *libvirt.Libvirt) string {
	v, err := l.ConnectGetLibVersion()
	if err != nil {
		log.Fatalf("Failed to retrieve libvirt version : %v", err)
	}
	//return string(v)
	return fmt.Sprint(v)
}

func getDomains(l *libvirt.Libvirt, f string) []libvirt.Domain {
	//domains, err := l.Domains()
	var flags libvirt.ConnectListAllDomainsFlags
	if f == "active" {
		flags = libvirt.ConnectListDomainsActive
	} else if f == "inactive" {
		flags = libvirt.ConnectListDomainsInactive
	} else {
		flags = libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive
	}

	domains, _, err := l.ConnectListAllDomains(1, flags)
	if err != nil {
		log.Fatalf("failed to retrieve domains: %v", err)
	}

	return domains
} */
