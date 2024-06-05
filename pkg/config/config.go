package config

import (
	"os"

	"github.com/DuvanAlbarracin/movies_profiles/pkg/utils"
)

type Config struct {
	Port  string
	DBUrl string
}

func LoadConfig() (c Config, err error) {
	portData, err := os.ReadFile("/run/secrets/profiles_port")
	if err != nil {
		return
	}

	dbUrlData, err := os.ReadFile("/run/secrets/db_url")
	if err != nil {
		return
	}

	c.Port = utils.TrimString(string(portData))
	c.DBUrl = utils.TrimString(string(dbUrlData))
	return
}
