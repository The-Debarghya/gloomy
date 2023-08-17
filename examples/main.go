package main

import (
	"errors"
	"flag"
	"os"

	"github.com/The-Debarghya/gloomy"
)

const logPath = "./example.log"

var verbose = flag.Bool("verbose", true, "print info level logs to stdout")

func doSomething() error {
	return errors.New("test error")
}

func main() {
	flag.Parse()

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		gloomy.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()

	defer gloomy.Init("GloomyExample", *verbose, true, lf).Close()

	gloomy.Info("I'm about to do something!")
	if err := doSomething(); err != nil {
		gloomy.Errorf("Error running doSomething: %v", err)
	}
}
