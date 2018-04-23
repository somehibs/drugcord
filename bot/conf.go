package drugcord

import (
	_"fmt"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-config/source/envvar"
	"github.com/micro/go-config/source/flag"
)

type PrometheusConfig struct {
	Host string
	Port int
	Enabled bool
}

// Need an initial request, auth by
type DiscordServer struct {
	InitialRequest string
	AuthenticatedBy string
}

type BotConfig struct {
	ID string `json:"id"`
	Token string
	Email string
	Password string
	Prometheus PrometheusConfig
	DiscordServers []DiscordServer `json:"discord_servers"`
}

func GetConf() BotConfig {
	var conf = config.NewConfig()
	conf.Load(envvar.NewSource(), flag.NewSource(), file.NewSource(file.WithPath("./config.json")))
	var botConf = BotConfig{}
	conf.Get().Scan(&botConf)
	return botConf
}
