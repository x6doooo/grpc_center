package env

import (
    "flag"
    "github.com/BurntSushi/toml"
)

const (
    ENV_MODE_DEV = "dev"
    ENV_MODE_PROD = "release"
)

var (
    IsDev bool
    IsProd bool
    confFile string
    Conf = &MainConf{}
)

func init() {
    flag.StringVar(&confFile, "conf", "/path/file.toml", "config file")
    flag.Parse()

    if confFile == "" {
        panic("need args --conf=/path/file")
    }
    toml.DecodeFile(confFile, Conf)
    IsDev = Conf.Env.Mode == ENV_MODE_DEV
    IsProd = Conf.Env.Mode == ENV_MODE_PROD
}

