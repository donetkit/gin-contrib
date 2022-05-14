package server

import (
	"os"
	"strings"
)

const (
	Dev  = "dev"
	Prod = "prod"
	Test = "test"
)

var EnvName = getEnvName()

func getEnvName() string {
	name := Dev
	switch strings.ToLower(os.Getenv("ENV")) {
	case Dev:
		name = Dev
	case Test:
		name = Test
	case Prod:
		name = Prod
	}
	return name
}

func IsDevelopment() bool {
	return EnvName == Dev
}

func IsTest() bool {
	return EnvName == Test
}

func IsProduction() bool {
	return EnvName == Prod
}
