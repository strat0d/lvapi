module github.com/strat0d/lvapi

go 1.13

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/strat0d/lvapi/lvget v0.0.0-00010101000000-000000000000
)

replace github.com/strat0d/lvapi/lvstr => ./lvstr

replace github.com/strat0d/lvapi/lvget => ./lvget
