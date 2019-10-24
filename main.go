package main

import (
	"fmt"
	"os"
)

func run() (string, error) {
	return "PROMPT='%n@%m %1~ %# '", nil
}

func main() {
	s, err := run()
	if err != nil {
		fmt.Print("PROMPT='failed > '")
		os.Exit(1)
	}
	fmt.Print(s)
}
