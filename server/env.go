package server

import "os"

const (
	Dev  = "dev"
	Prod = "prod"
	Test = "test"
)

var EnvName = getEnvName()

func getEnvName() string {
	name := os.Getenv("ENV")
	if name == "" {
		return Dev
	}
	return name
}
