package config

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

var (
	Database       *sql.DB
	countConnected uint = 1
	MAX_CONNECTION uint = 5
)

type config_DB struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type config struct {
	DB config_DB `yaml:"db"`
}

func (c *config_DB) createDatabaseString() string {
	var result strings.Builder

	result.WriteString(c.Username)
	result.WriteString(":" + c.Password)
	result.WriteString("@tcp(" + c.Host)
	if c.Port != "" {
		result.WriteString(":" + c.Port)
	}
	result.WriteString(")/" + c.DBName)

	return result.String()
}

func Connect() {
	bitesInfo, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		panic(err)
	}

	var conf config

	if err = yaml.Unmarshal(bitesInfo, &conf); err != nil {
		panic(err)
	}

	Database, err = sql.Open("mysql", conf.DB.createDatabaseString())

	var (
		isError           = err != nil
		isLimitConnection = countConnected > MAX_CONNECTION
	)

	switch {
	case isError && !isLimitConnection:
		countConnected += 1
		Connect()
		break
	case isError:
		panic(err)
	case !isError:
		countConnected = 0
		log.Print("Database was connected")
	}
}

func Disconnect() {
	if Database.Ping() == nil {
		return
	}

	if err := Database.Close(); err != nil {
		panic(err)
	} else {
		log.Print("Database was disconnected")
	}
}
