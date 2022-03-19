package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/strat0d/lvapi/lvget"
	"libvirt.org/go/libvirt"
)

type Host struct {
	driver string
	user   string
	host   string
	level  string
}

func defaultHost(host string) *Host {
	newHost := Host{driver: "qemu+ssh", user: "root", level: "system"}
	newHost.host = host
	return &newHost
}

func (h Host) URI() string {
	return fmt.Sprintf("%s://%s@%s/%s", h.driver, h.user, h.host, h.level)
}

func lvapiConn(host string, write bool) (*libvirt.Connect, error) {
	h := defaultHost(host)
	var f libvirt.ConnectFlags
	if !write {
		f = libvirt.CONNECT_RO
	}
	conn, err := libvirt.NewConnectWithAuthDefault(h.URI(), libvirt.ConnectFlags(f))
	return conn, err
}

func addHostToMap(h string, m map[string]*libvirt.Connect) error {
	if m[h] != nil {
		// already exists
		if alive, _ := m[h].IsAlive(); alive {
			// connection still alive
			return nil
		}
	}
	conn, err := lvapiConn(h, true)
	if err != nil {
		//log.Fatalf("Failed to open libvirt to %v: %v", h, err)
		return err
	}
	log.Printf("Opened libvirt to %v", h)
	m[h] = conn
	return nil
}

func main() {
	//GIN
	router := gin.Default()
	router.SetTrustedProxies(nil)

	lvcs := make(map[string]*libvirt.Connect)

	ag_domains := router.Group("/api/v0/domains")
	{
		//get all domains on a host
		ag_domains.GET("/:host", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")

			res := make(chan lvget.DomainsResult)
			h := c.Param("host")
			if err := addHostToMap(h, lvcs); err != nil {
				c.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				go func() {
					res <- lvget.Domains(lvcs[h])
					close(res)
				}()
				r := <-res
				if r.Error != nil {
					c.IndentedJSON(http.StatusOK, gin.H{"error": r.Error.Error()})
				} else {
					c.IndentedJSON(http.StatusOK, r.Domains)
				}
			}

		})
		//get domain :val :by(id, name, uuid) on :host
		ag_domains.GET("/:host/:by/:val", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")

			res := make(chan lvget.DomainResult)
			h, by, v := c.Param("host"), c.Param("by"), c.Param("val")
			if err := addHostToMap(h, lvcs); err != nil {
				c.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				go func() {
					res <- lvget.Domain(lvcs[h], by, v)
					close(res)
				}()
				r := <-res
				if r.Error != nil {
					c.IndentedJSON(http.StatusOK, gin.H{"error": r.Error.Error()})
				} else {
					c.IndentedJSON(http.StatusOK, r.Domain)
				}
			}
		})

		//POST
		//run libvirt/virsh :action on guest :val found :by<id,name,uuid> on :host
		ag_domains.POST("/:host/:by/:val/:action", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")

			dom := make(chan lvget.LvDomainResult)
			h, by, v, a := c.Param("host"), c.Param("by"), c.Param("val"), c.Param("action")
			if err := addHostToMap(h, lvcs); err != nil {
				c.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				go func() {
					dom <- lvget.LvDomain(lvcs[h], by, v)
					close(dom)
				}()
				d := <-dom
				if d.Error != nil {
					c.IndentedJSON(http.StatusOK, gin.H{"error": d.Error.Error()})
				} else {
					switch a {
					case "destroy":
						c.IndentedJSON(http.StatusOK, d.Domain.Destroy())
					case "reboot":
						c.IndentedJSON(http.StatusOK, d.Domain.Reboot(0))
					case "reset":
						c.IndentedJSON(http.StatusOK, d.Domain.Reset(0))
					case "resume":
						c.IndentedJSON(http.StatusOK, d.Domain.Resume())
					case "start":
						c.IndentedJSON(http.StatusOK, d.Domain.Create())
					case "suspend":
						c.IndentedJSON(http.StatusOK, d.Domain.Suspend())
					case "shutdown":
						c.IndentedJSON(http.StatusOK, d.Domain.Shutdown())
					default:
						c.IndentedJSON(http.StatusOK, gin.H{"error": fmt.Sprintf("%v not implemented", a)})
					}
					c.IndentedJSON(http.StatusOK, d.Domain.Destroy())
				}
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
	for _, h := range lvcs {
		h.Close()
	}
}
