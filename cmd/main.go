package main

import (
	"flag"
	"os"
	"os/signal"

	solver "github.com/xdqc/letterpress-solver"
)

var ()

func init() {
	flag.Parse()
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		solver.Run("8998")
	}()
	go func() {
		solver.RunWeb("8080")
	}()
	<-c
	solver.Close()
	return
}
