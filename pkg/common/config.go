package common

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/hashicorp/consul/api"
)

// 保持你的全局变量定义
var (
	flagConf   string
	consulAddr string
	consulPath string
)

func init() {
	flag.StringVar(&flagConf, "conf", "../../configs", "local config path")
	flag.StringVar(&consulAddr, "consul_addr", "", "consul address")
	flag.StringVar(&consulPath, "consul_path", "", "consul config path")
}

func BootstrapConfig(serviceName string) config.Config {
	// 优先级：Flag > Env
	addr := consulAddr
	if addr == "" {
		addr = os.Getenv("CONSUL_ADDR")
	}

	path := consulPath
	if path == "" {
		path = os.Getenv("CONSUL_PATH")
		if path == "" {
			path = "configs/" + serviceName
		}
	}

	var sources []config.Source

	// 永远携带本地文件源，作为默认兜底
	sources = append(sources, file.NewSource(flagConf))

	if addr != "" {
		c := api.DefaultConfig()
		c.Address = addr
		client, err := api.NewClient(c)
		if err != nil {
			panic(err)
		}
		source, err := consul.New(client, consul.WithPath(path))
		if err != nil {
			panic(err)
		}
		sources = append(sources, source)
	}

	c := config.New(config.WithSource(sources...))

	if err := c.Load(); err != nil {
		panic(err)
	}

	return c
}
