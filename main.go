package main

import (
	"log"
	"tax-calculator-cli/internal/cli"
)

func main() {
	processor := cli.NewProcessor()
	if err := processor.Run(); err != nil {
		log.Fatal(err)
	}
}
