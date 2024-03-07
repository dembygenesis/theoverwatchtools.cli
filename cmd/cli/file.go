package main

import (
	"github.com/dembygenesis/local.tools/internal/lib/fslib"
	"github.com/spf13/cobra"
)

var copyFolderAToBCommand = &cobra.Command{
	Use:   copyFolderAToB.string(),
	Short: "Copies files from folder 'A', to folder 'B'.",
	Long: `
        This command streamlines the process of transferring files from a source directory (Folder 'A') to a target directory (Folder 'B'), while providing a versatile set of options for customized operation. The functionalities include:

        1. **Selective Copying with GenericExclusions from Folder A:**
           - Facilitates selective copying from the source directory, allowing for specific exclusions. Users can define files or subdirectories in Folder 'A' that should be omitted from the copying process, ensuring that only pertinent files are included in the operation.

        2. **Pre-Copy Cleanup with WipeFolderB:**
           - Offers an option to perform a cleanup of the destination directory (Folder 'B') prior to copying. If enabled, the command will clear all contents of Folder 'B', paving the way for a clean slate that will exclusively contain the files transferred from Folder 'A'.

        The command initiates with a preface operation, ensuring all conditions are met for a smooth and error-free file transfer. Post this preliminary step, the command meticulously logs each phase of the operation, ensuring transparency and traceability of the process flow.
    `,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		dst := args[1]

		opts := fslib.CopyOptions{
			Source:      args[0],
			Destination: args[1],
		}

		if err := srv.CopyDirToAnother(&opts); err != nil {
			log.Errorf("clip gpt preface: %v", err)
			return
		}

		log.Infof("Copied '\033[1m%s\033[0m' to '\033[1m%s\033[0m'", src, dst)
	},
}
