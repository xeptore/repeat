package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"sync"

	"golang.org/x/sync/semaphore"
)

func executeCommand(shell, command string) {
	cmd := exec.Command(shell, "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func makeRepeater(conf *Config) func() {
	return func() {
		sem := semaphore.NewWeighted(int64(conf.Parallelism))
		ctx := context.Background()
		wg := sync.WaitGroup{}

		var i uint
		for i = 0; i < conf.Iterations; i++ {
			wg.Add(1)
			sem.Acquire(ctx, 1)
			go func(i uint) {
				defer sem.Release(1)
				defer wg.Done()
				executeCommand(conf.Shell, conf.Command)
			}(i)
		}

		wg.Wait()
	}
}

func main() {
	conf, err := loadConfig()
	if nil != err {
		log.Fatalln(err.Error())
	}

	repeat := makeRepeater(&conf)
	repeat()
}
