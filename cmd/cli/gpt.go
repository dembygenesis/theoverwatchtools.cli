package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var copyGptCodePrefaceToClipboardCommand = &cobra.Command{
	Use:   clipGptPreface,
	Short: "Copies a code preface for chat GPT that ensures code quality.",
	Long: `
		Copies a code preface for chat GPT that ensures code quality.
		ChatGPT usually returns "decent code" if it is on ChatGPT 4, but the good engineering
		foundations usually still has something to be desired.

		It lacks:
		Defensive programming, testability, readability, modularity - and this
		preface attempts to remediate that. It obviously will not be perfect,
		but it gives tangible improvements (at least based on anecdotal experience).
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := srv.ClipCodingStandardsPreface(); err != nil {
			return fmt.Errorf("clip gpt preface: %v", err)
		}
		logger.Info("Copied gpt preface")

		return nil
	},
}
