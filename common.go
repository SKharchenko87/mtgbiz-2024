package common

import (
	"fmt"
	"os"
	"time"
)
import _ "encoding/json"

type Table1 struct {
	Id         int64     `json:"id"`
	N          int64     `json:"n"`
	Code       [4]byte   `json:"code"`
	Data       string    `json:"data"`
	CreateDttm time.Time `json:"createDttm"`
}

func GetDSN() string {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	return dsn
}

func GetServerURL(postfix string) string {
	protocol := os.Getenv("SERVER_PROTOCOL_" + postfix)
	host := os.Getenv("SERVER_HOST_" + postfix)
	port := os.Getenv("SERVER_PORT_" + postfix)
	pathParam := os.Getenv("SERVER_PATH_PARAM_" + postfix)
	serverURL := fmt.Sprintf("%s://%s:%s/%s", protocol, host, port, pathParam)
	return serverURL
}
