package postgres

import (
	"encoding/json"
	"os"
)
type dbConfig struct {
	Addr     string `json:"host"`
	Port     int	`json:"port"`
	Username string	`json:"username"`
	Password string	`json:"password"`
	DBName     string `json:"dbname"`
}
func getDbConfig() *dbConfig {
	config := dbConfig{}
	file := "./configs/config.json"
	data, err := os.ReadFile(file)
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
