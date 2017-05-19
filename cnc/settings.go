package main

import (
	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego/logs"
)

type SpringSettings struct {
	Title         string
	Version       string
	Server        string
	ServerCrt     string
	ServerKey     string
	AdminServer   string
	AdminName     string
	AdminPassword string
	Log           LogSettings
}

func NewSettings(settingsFile string, settings *SpringSettings) error {
	_, err := toml.DecodeFile(settingsFile, settings)
	return err
}

type LogSettings struct {
	Stdout bool
	Path   string
	Level  string `toml:"Level"`
}

var LogLevelMap = map[string]int{
	"DEBUG":  logs.LevelDebug,
	"INFO":   logs.LevelInfo,
	"NOTICE": logs.LevelNotice,
	"WARN":   logs.LevelWarning,
	"ERROR":  logs.LevelError,
}

func (ls LogSettings) BeeLevel() int {
	l, ok := LogLevelMap[ls.Level]
	if !ok {
		panic("Config error: invalid log level: " + ls.Level)
	}
	return l
}
