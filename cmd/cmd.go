package cmd

import (
	"flag"
	"fmt"
	"github.com/sergeyignatov/simpleipam/api"
	icfg "github.com/sergeyignatov/simpleipam/config"
	ctx "github.com/sergeyignatov/simpleipam/context"
	"github.com/sergeyignatov/simpleipam/subnet"
	"net/http"
	"os"
	"path"
)

func Run() error {

	var err error

	listen := flag.String("listen", ":4567", "host:port for HTTP listening")
	cfg := flag.String("cfg", "config.yml", "config")
	flag.Parse()

	config, err := icfg.LoadConfig(*cfg)
	if err != nil {
		return err
	}
	if _, err := os.Stat(config.DataDir); err != nil {
		return fmt.Errorf("Directory %s is not exists", config.DataDir)
	}
	f, err := os.Create(path.Join(config.DataDir, ".check_perm"))
	if err != nil {
		return fmt.Errorf("check write permissions in %s", config.DataDir)
	}
	os.Remove(path.Join(config.DataDir, ".check_perm"))
	defer f.Close()
	subnets := subnet.NewSubnets()
	subnets.Load(config)
	context := ctx.Context{config, subnets}
	err = http.ListenAndServe(*listen, api.Router(&context))
	return err
}
