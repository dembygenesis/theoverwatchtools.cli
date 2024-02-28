package main

import (
	"github.com/dembygenesis/local.tools/di/cfg"
	"github.com/sarulabs/dingo/v4"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	err := dingo.GenerateContainer((*cfg.Provider)(nil), os.Args[1])
	if err != nil {
		log.Println("error generating container: ", err)
		os.Exit(1)
	}
}
