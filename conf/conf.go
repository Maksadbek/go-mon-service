package conf

import (
	"github.com/BurntSushi/toml"
	"io"
)

type Datastore struct {
	Redis struct {
		Port string
		Chan string
	}
}

type Server struct {
	IP   string
	Port string
}

type App struct {
	DS  Datastore
	SRV Server
	Log Log
}

type Log struct {
	Attrs []string
}

func Read(r io.Reader) (config App, err error) {
	_, err = toml.DecodeReader(r, &config)
	return config, err
}