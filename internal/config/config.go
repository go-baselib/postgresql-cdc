package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"time"
)

// Config contains ENV variables
type Config struct {
	Source     `yaml:"source"`
	EventQueen `yaml:"event_queen"`
}

type Source struct {
	DbHost string `yaml:"db_host"`
	DbPort string `yaml:"db_port"`
	DbUser string `yaml:"db_user"`
	DbPass string `yaml:"db_pass"`
	DbName string `yaml:"db_name"`

	Slot string `yaml:"slot"`
	Lsn  string `yaml:"lsn"`

	TableNames string `yaml:"table_names"`
	Chunks     string `yaml:"chunks"`
}

type EventQueen struct {
	Type          string        `yaml:"type"`
	Hosts         string        `yaml:"hosts"`
	Name          string        `yaml:"name"`
	MaxRetry      int           `yaml:"max_retry"`
	RetryInterval time.Duration `yaml:"retry_interval"` // 单位ms
}

var gCfg *Config

func InitConfig(in []byte) {
	var c Config
	if err := yaml.Unmarshal(in, &c); err != nil {
		panic(err)
	}

	if c.Slot == "" {
		panic("slot is not set")
	}

	if c.RetryInterval <= 0 {
		c.RetryInterval = 3000
	}

	if c.Lsn == "" {
		c.Lsn = "0/0"
	}

	if c.MaxRetry <= 0 {
		c.MaxRetry = 3
	}

	log.Printf("config: %+v\n", c)
	gCfg = &c
}

func GetConfig() *Config {
	return gCfg
}
