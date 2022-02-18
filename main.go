package main

import (
	"github.com/strat0d/lvapi/lvxml"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"libvirt.org/go/libvirt"
)

func main() {
	//GIN
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.StaticFile("/ep.html", "./static/ep.html")
	router.GET("/defaultxml", func(c *gin.Context) {
		var dom lvxml.Domain
		lvxml.GetDefaultDomainXML(&dom)
		c.XML(http.StatusOK, dom)
	})
	router.GET("/domains", func(c *gin.Context) {
		lvcf := libvirt.ConnectFlags(libvirt.CONNECT_RO)
		lvconn, err := libvirt.NewConnectWithAuthDefault("qemu+ssh://strat@192.168.101.2/system", lvcf)
		if err != nil {
			log.Fatalf("Error connecting to libvirt: %v", err)
		}
		defer lvconn.Close()

		domains, err := lvconn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
		if err != nil {
			log.Fatalf("Error getting domains: %v", err)
		}

		type domS struct {
			DomainName string
			DomainStatus int
			DomainID uint
		}
		doms := []domS{}
		for _, dom := range domains {
			name, err := dom.GetName()
			if err != nil {
				log.Fatalf("err: %v", err)
			}
			_, status, err := dom.GetState()
			if err != nil {
				log.Fatalf("err: %v", err)
			}
			var id uint = 0
			if status > 0 {
				var err error
				id, err = dom.GetID()
				if err != nil {
					log.Fatalf("err: %v", err)
				}
			} else {
				id = 0
			}
			d := domS{DomainName: name, DomainStatus: status, DomainID: id}
			doms = append(doms, d)
		}
		c.IndentedJSON(http.StatusOK, doms)
	})

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
