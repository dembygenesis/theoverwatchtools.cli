package main

import (
	"github.com/spf13/cobra"
)

var copyToClipboardCmd = &cobra.Command{
	Use:   clipFileContents.string(),
	Short: "Copy to clipboard copies from the directory provided.",
	Long:  `Copies files contents from the root path provided.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		files, err := srv.CopyToClipboard(args[0], nil)
		if err != nil {
			logger.Errorf("copy to clipboard: %v", err)
			return
		}
		logger.Infof("copied \033[1;34m%v\033[0m files to clipboard!", len(files))
	},
}
