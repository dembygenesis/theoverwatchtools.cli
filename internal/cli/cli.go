package cli

import (
	"github.com/dembygenesis/local.tools/di/ctn/dic"
	"log"
)

type Cli struct {
	ctn *dic.Container

	// Also has cobra, so where do I execute it? I merely execute it in the cmd, but
	// I do my functionalities here
}

func New() (*Cli, error) {
	// Load container
	ctn, err := dic.NewContainer()
	if err != nil {
		log.Fatalf("new container: %v", err)
	}

	// Load commands
	cli := &Cli{ctn: ctn}

	// Where the fuck do I load hue hue

	return cli, nil
}
