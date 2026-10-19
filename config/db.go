package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Read config.yaml from the application's working directory.
	v := viper.New()
	v.SetConfigFile("config.yaml")

	if err := v.ReadInConfig(); err != nil {
		panic("Unable to read config.yaml: " + err.Error())
	}

	host := v.GetString("mysql.host")
	port := v.GetInt("mysql.port")
	username := v.GetString("mysql.username")
	password := v.GetString("mysql.password")
	database := v.GetString("mysql.database")

	// Check the required settings before connecting.
	if strings.TrimSpace(host) == "" ||
		strings.TrimSpace(username) == "" ||
		password == "" ||
		strings.TrimSpace(database) == "" {
		panic("MySQL host, username, password and database are required")
	}

	if port < 1 || port > 65535 {
		panic("MySQL port must be between 1 and 65535")
	}

	// Let the MySQL driver format the connection string.
	dbConfig := mysqldriver.NewConfig()
	dbConfig.User = username
	dbConfig.Passwd = password
	dbConfig.Net = "tcp"
	dbConfig.Addr = net.JoinHostPort(host, strconv.Itoa(port))
	dbConfig.DBName = database
	dbConfig.ParseTime = true
	dbConfig.Params = map[string]string{
		"charset": "utf8mb4",
	}

	db, err := gorm.Open(mysql.Open(dbConfig.FormatDSN()), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic(fmt.Errorf("database connection failed: %w", err))
	}

	DB = db
}
