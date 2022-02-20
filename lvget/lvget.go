package lvget

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/strat0d/lvapi/lvstr"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func DefaultXML() *libvirtxml.Domain {
	dom := &libvirtxml.Domain{Type: "kvm", Name: "TestName"}
	return dom
}

func defaultHost(host string) *Host {
	newHost := Host{driver: "qemu+ssh", user: "root", level: "system"}
	newHost.host = host
	return &newHost
}

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

type DomainsResult struct {
	Domains []lvstr.Domain
	Err     error
}

func Domains(c *gin.Context) DomainsResult {
	lvconn, err := lvapiConnRo(c.Param("host"))
	if err != nil {
		return DomainsResult{Domains: nil, Err: err}
		//c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("%v", err)})
	}
	defer lvconn.Close()

	domains, err := lvconn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		return DomainsResult{Domains: nil, Err: err}
		//c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("Error getting domains: %v", err)})
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

	//c.IndentedJSON(http.StatusOK, doms)
	return DomainsResult{Domains: doms, Err: nil}
}

type DomainResult struct {
	Domain lvstr.Domain
	Err    error
}

func Domain(c *gin.Context) DomainResult {
	h := c.Param("host")
	m := c.Param("method")
	v := c.Param("val")

	lvconn, err := lvapiConnRo(h)
	if err != nil {
		//c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): %v", err)})
		return DomainResult{Err: err}
	}
	defer lvconn.Close()

	var d = lvstr.Domain{}

	switch m {
	case "id":
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			//c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): \"%v\" invalid ID. %v", v, err)})
			return DomainResult{Err: err}
		}
		dom, err := lvconn.LookupDomainById(uint32(id))
		if err != nil {
			//c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): \"%v\"", err)})
			return DomainResult{Err: err}
		}
		lvstr.GetDomain(dom, &d)
	case "name":
		//
		dom, err := lvconn.LookupDomainByName(v)
		if err != nil {
			//c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): \"%v\"", err)})
			return DomainResult{Err: err}
		}
		lvstr.GetDomain(dom, &d)
	case "uuid":
		//
		dom, err := lvconn.LookupDomainByUUIDString(v)
		if err != nil {
			//c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): \"%v\"", err)})
			return DomainResult{Err: err}
		}
		lvstr.GetDomain(dom, &d)
	default:
		//c.JSON(http.StatusOK, gin.H{"error": fmt.Sprintf("getDomain(): \"%v\" invalid method", m)})
		return DomainResult{Err: errors.New("invalid method")}
	}
	//c.IndentedJSON(http.StatusOK, d)
	return DomainResult{d, nil}
}
