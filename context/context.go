package context

import (
	"github.com/sergeyignatov/simpleipam/config"
	"github.com/sergeyignatov/simpleipam/subnet"
)

type Context struct {
	Config  *config.Config
	Subnets *subnet.Subnets
}
