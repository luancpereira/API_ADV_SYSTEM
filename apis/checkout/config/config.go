package config

import "os"

var (
	SERVER_PORT         = os.Getenv("SERVER_PORT")
	SWAGGER_SERVER_HOST = os.Getenv("SWAGGER_SERVER_HOST")
)
