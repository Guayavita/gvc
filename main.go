package main

import (
	"github.com/charmbracelet/log"
	"jmpeax.com/guayavita/gvc/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Errorf("%v", err.Error())
	}
}
