package flags

import (
	"flag"

	"github.com/caarlos0/env/v6"

	"github.com/ramil063/firstgodiplom/internal/logger"
)

var RunAddress = "localhost:8080"
var AccrualSystemAddress = "http://localhost:8081"
var DatabaseURI = ""

type EnvVars struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
}

func ParseFlags() {
	flag.StringVar(&RunAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&DatabaseURI, "d", "", "database URI")
	flag.StringVar(&AccrualSystemAddress, "r", "http://localhost:8081", "address and port to run accrual server")
	flag.Parse()

	var ev EnvVars

	_ = env.Parse(&ev)

	if ev.RunAddress != "" {
		RunAddress = ev.RunAddress
	}

	if ev.AccrualSystemAddress != "" {
		AccrualSystemAddress = ev.AccrualSystemAddress
	}

	if ev.DatabaseURI != "" {
		DatabaseURI = ev.DatabaseURI
	}

	logger.WriteInfoLog("RunAddress:" + RunAddress)
	logger.WriteInfoLog("AccrualSystemAddress:" + AccrualSystemAddress)
	logger.WriteInfoLog("DatabaseURI:" + DatabaseURI)
}
