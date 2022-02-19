module github.com/strat0d/lvapi

go 1.13

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/strat0d/lvapi/lvstr v0.0.0-20220219201353-4df22a246d41
	libvirt.org/go/libvirt v1.7010.0
	libvirt.org/go/libvirtxml v1.7010.0
)

replace github.com/strat0d/lvapi/lvstr => ./lvstr
