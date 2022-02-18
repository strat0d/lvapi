package main

import (
	"fmt"
	"log"

	"lvxml"

	"github.com/gin-gonic/gin"
	"libvirt.org/go/libvirt"
)

func main() {
	//GIN
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.StaticFile("/ep.html", "./static/ep.html")

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

	for _, dom := range domains {
		name, err := dom.GetName()
		status, err := dom.GetInfo()
		if err == nil {
			fmt.Printf("%s - %d\n", name, status.State)
		}

	}
	dom := lvxml.GetDefaultDomainXML()
	fmt.Printf("%v", dom)

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
