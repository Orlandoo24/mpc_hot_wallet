package config

import "github.com/zeromicro/go-zero/rest"

type ChainConf struct {
	Name    string `json:"Name"`
	RpcUrl  string `json:"RpcUrl"`
	ChainId int64  `json:"ChainId"`
}

type Config struct {
	rest.RestConf
	Postgres struct {
		DSN string
	}
	Lifi struct {
		ApiUrl string
	}
	// Chains maps a chain name (e.g., "BSC") to its configuration.
	Chains map[string]ChainConf
}
