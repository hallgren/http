package main

import (
	"fmt"
	"os"

	"github.com/lunjon/http/command"
)

func main() {
	err := command.Build("0.10.0").Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
