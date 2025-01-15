package config

import (
	"os"
)

var (
	ERROR_FILE = os.Getenv("ERROR_FILE")
)
