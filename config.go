package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/alexflint/go-arg"
	"github.com/xeptore/gologger"
)

type Config struct {
	Parallelism uint   `arg:"env:PARALLELISM"`
	Command     string `arg:"positional"`
	Shell       string `arg:"env:SHELL"`
	Iterations  uint   `arg:"env:ITERATIONS"`
}

func getDefaultParallelism() uint {
	return uint(runtime.NumCPU())
}

func getParallelism(conf *Config) uint {
	if parallelism := conf.Parallelism; parallelism == 0 {
		defaultParallelism := getDefaultParallelism()
		gologger.Warn(fmt.Sprintf("Parallelism amount is not specified, default value of %d will be used", defaultParallelism))
		return defaultParallelism
	}

	if parallelism := conf.Parallelism; parallelism > uint(runtime.NumCPU()) {
		gologger.Warn(fmt.Sprintf("Parallelism is set to %d, while number of machine CPU cores are %d", parallelism, runtime.NumCPU()))
	}

	return conf.Parallelism
}

func ensureWorkingShell(conf *Config) {
	cmd := exec.Command(conf.Shell, "-c", "echo")
	err := cmd.Run()
	if nil != err {
		gologger.ErrorFatal(fmt.Sprintf("Could not execute a test shell script via '-c' option to provided $SHELL"))
	}
}

func loadConfig() (Config, error) {
	conf := Config{}

	arg.MustParse(&conf)

	conf.Parallelism = getParallelism(&conf)

	ensureWorkingShell(&conf)

	return conf, nil
}
