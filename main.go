package main

import (
	"jmpeax.com/guayavita/gvc/cmd"
	"jmpeax.com/guayavita/gvc/internal/term"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		term.Error(err.Error())
	}
}
