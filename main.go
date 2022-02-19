package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"

	"github.com/strat0d/lvapi/lvstr"
	//"lvapi/lvstr"
)

func lvapiConnRo(host string) (*libvirt.Connect, error) {
	h := defaultHost(host)
	conn, err := libvirt.NewConnectWithAuthDefault(h.URI(), libvirt.ConnectFlags(libvirt.CONNECT_RO))
	return conn, err
}

type Host struct {
	driver string
	user   string
	host   string
	level  string
}

func (h Host) URI() string {
	return fmt.Sprintf("%s://%s@%s/%s", h.driver, h.user, h.host, h.level)
}

func defaultHost(host string) *Host {
	newHost := Host{driver: "qemu+ssh", user: "root", level: "system"}
	newHost.host = host
	return &newHost
}

func getDefaultXML() *libvirtxml.Domain {
	dom := &libvirtxml.Domain{Type: "kvm", Name: "TestName"}
	return dom
}

func getDomains(c *gin.Context) {
	lvconn, err := lvapiConnRo(c.Param("host"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("%v", err)})
	}
	defer lvconn.Close()

	domains, err := lvconn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("Error getting domains: %v", err)})
	}

	var wg sync.WaitGroup

	doms := []lvstr.Domain{}
	wg.Add(len(domains))
	for _, dom := range domains {
		d := lvstr.Domain{}
		go func(domGr libvirt.Domain) {
			defer wg.Done()
			lvstr.GetDomain(&domGr, &d)
			doms = append(doms, d)
		}(dom)
	}
	wg.Wait()

	//always return domaisn in alphabetical order by Name
	sort.Slice(doms, func(i, j int) bool {
		return doms[i].Name < doms[j].Name
	})
	c.IndentedJSON(http.StatusOK, doms)
}

func getDomain(c *gin.Context) {
	h := c.Param("host")
	m := c.Param("method")
	v := c.Param("val")

	lvconn, err := lvapiConnRo(h)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): %v", err)})
	}
	defer lvconn.Close()

	var d = lvstr.Domain{}

	switch m {
	case "id":
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): \"%v\" invalid ID. %v", v, err)})
			return
		}
		dom, err := lvconn.LookupDomainById(uint32(id))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): \"%v\"", err)})
			return
		}
		lvstr.GetDomain(dom, &d)
	case "name":
		//
		dom, err := lvconn.LookupDomainByName(v)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): \"%v\"", err)})
			return
		}
		lvstr.GetDomain(dom, &d)
	case "uuid":
		//
		dom, err := lvconn.LookupDomainByUUIDString(v)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): \"%v\"", err)})
			return
		}
		lvstr.GetDomain(dom, &d)
	default:
		c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): \"%v\" invalid method", m)})
		return
	}
	c.IndentedJSON(http.StatusOK, d)
}

func main() {
	//GIN
	router := gin.Default()
	//router.SetTrustedProxies(nil)

	ag_domains := router.Group("/api/v0/domains")
	{
		//get all domains on a host
		ag_domains.GET("/:host", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			getDomains(c)
		})
		//get domain :val by :<method>(id, name, uuid) on :host
		ag_domains.GET("/:host/:method/:val", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			getDomain(c)
		})

		//POST
		ag_domains.POST("/:host/", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			//
		})
	}

	ag_misc := router.Group("/api/v0/misc")
	{
		ag_misc.GET("/defaultxml", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			dom := getDefaultXML()
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
