package main

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/di/ctn/dic"
	"github.com/dembygenesis/local.tools/internal/cli"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/spf13/cobra"
	"os"
)

var (
	ctn *dic.Container
	err error
	srv *cli.Service
	log = logger.New(context.TODO())
)

func init() {
	ctn, err = dic.NewContainer()
	if err != nil {
		log.Fatalf("new container: %v", err)
	}

	srv, err = ctn.SafeGetServicesLayer()
	if err != nil {
		log.Fatalf("get services: %v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "Root",
	Short: "This is the root command.",
	Long:  `This is the root command.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Please select a command")
	},
}

func init() {
	rootCmd.AddCommand(copyToClipboardCmd)
	rootCmd.AddCommand(copyGptCodePrefaceToClipboardCommand)
	rootCmd.AddCommand(copyFolderAToBCommand)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		os.Exit(1)
	}
}
