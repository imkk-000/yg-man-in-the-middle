package model

import (
	"encoding/json"
	"io/ioutil"
)

type IpConfig struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type Config struct {
	Server IpConfig `json:"server"`
	Login  IpConfig `json:"login"`
	Game   IpConfig `json:"game"`
	Log    bool     `json:"log"`
}

func (Config) Read(filename string) (config Config) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return
}
