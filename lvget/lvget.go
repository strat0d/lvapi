package lvget

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"

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

func lvapiConn(host string, write bool) (*libvirt.Connect, error) {
	h := defaultHost(host)
	var f libvirt.ConnectFlags
	if !write {
		f = libvirt.CONNECT_RO
	}
	conn, err := libvirt.NewConnectWithAuthDefault(h.URI(), libvirt.ConnectFlags(f))
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
	Error   error
}

func Domains(lvconn *libvirt.Connect) DomainsResult {
	//lvconn, err := lvapiConn(h, false)
	//if err != nil {
	//	return DomainsResult{Domains: nil, Error: err}
	//}
	// defer lvconn.Close()

	domains, err := lvconn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		return DomainsResult{Domains: nil, Error: err}
	}
	if len(domains) == 0 {
		return DomainsResult{Domains: nil, Error: err}
	}

	// doms := make([]lvstr.Domain, 0, len(domains))
	// first_dom := domains[0]
	// d := lvstr.Domain{}
	// lvstr.GetDomain(&first_dom, &d)
	// doms = append(doms, d)

	// domains = append(domains[:0], domains[1:]...)

	var (
		mu   = &sync.Mutex{}
		doms = make([]lvstr.Domain, 0)
	)

	var wg sync.WaitGroup
	//wg.Add(len(domains))
	for _, dom := range domains {
		wg.Add(1)
		go func(domGr libvirt.Domain) {
			defer wg.Done()
			d := lvstr.Domain{}
			lvstr.GetDomain(&domGr, &d)
			mu.Lock()
			doms = append(doms, d)
			mu.Unlock()
		}(dom)
	}
	wg.Wait()
	/* 	for _, dom := range domains {
	   		d := lvstr.Domain{}
	   		lvstr.GetDomain(&dom, &d)
	   		doms = append(doms, d)
	   	}
	*/
	//always return domaisn in alphabetical order by Name
	sort.Slice(doms, func(i, j int) bool {
		return doms[i].Name < doms[j].Name
	})

	return DomainsResult{Domains: doms, Error: nil}
}

type DomainResult struct {
	Domain lvstr.Domain
	Error  error
}

type LvDomainResult struct {
	Domain *libvirt.Domain
	Error  error
}

func lvDomainById(c *libvirt.Connect, id uint32) (*libvirt.Domain, error) {
	d, err := c.LookupDomainById(id)
	if err != nil {
		return &libvirt.Domain{}, err
	}
	return d, nil
}
func lvDomainByName(c *libvirt.Connect, name string) (*libvirt.Domain, error) {
	d, err := c.LookupDomainByName(name)
	if err != nil {
		return &libvirt.Domain{}, err
	}
	return d, nil
}

func lvDomainByUUID(c *libvirt.Connect, uuid string) (*libvirt.Domain, error) {
	d, err := c.LookupDomainByUUIDString(uuid)
	if err != nil {
		return &libvirt.Domain{}, err
	}
	return d, nil
}

func LvDomain(lvconn *libvirt.Connect, by, v string) LvDomainResult {
	switch by {
	case "id":
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return LvDomainResult{nil, err}
		}
		d, err := lvDomainById(lvconn, uint32(id))
		if err != nil {
			return LvDomainResult{nil, err}
		}
		return LvDomainResult{d, nil}

	case "name":
		d, err := lvDomainByName(lvconn, v)
		if err != nil {
			return LvDomainResult{nil, err}
		}
		return LvDomainResult{d, nil}
	case "uuid":
		d, err := lvDomainByUUID(lvconn, v)
		if err != nil {
			return LvDomainResult{nil, err}
		}
		return LvDomainResult{d, nil}
	}
	return LvDomainResult{nil, errors.New("invalid method")}
}

func Domain(lvconn *libvirt.Connect, by, v string) DomainResult {
	ld := LvDomain(lvconn, by, v)
	if ld.Error != nil {
		return DomainResult{Error: ld.Error}
	}

	var d = lvstr.Domain{}
	lvstr.GetDomain(ld.Domain, &d)

	return DomainResult{d, nil}
}
