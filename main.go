package main

import (
	"fmt"
	"updater/internal/commands/update"
	"updater/internal/utils/error"
	"updater/internal/utils/filedeleter"
	"updater/internal/utils/system"
	"updater/internal/utils/unzipper"

	"github.com/inancgumus/screen"
	"github.com/spf13/cobra"
)


func main() {
    screen.Clear()
	screen.MoveTopLeft()
    
    var downloadPath, extractPath string

    rootCmd := &cobra.Command {
        Use: "vue-updater",
        Short: "A CLI tool to update vue ui torrent frontend",
        Run: func(cmd *cobra.Command, args []string) {
            downloadFilePath, err := update.RunCommand(downloadPath)
            if err != nil {
                error.Log(err, "Failed to download file")
            }

            if err := unzipper.UnzipWithProgress(downloadFilePath, extractPath); err != nil {
                error.Log(err, "Failed to unzip archive")
            }

            if err := filedeleter.DeleteFileWithProgress(downloadFilePath); err != nil {
                error.Log(err, "Failed to remove downloaded archive")
            }

            fmt.Println()
        },
    }

    rootCmd.Flags().StringVarP(&downloadPath, "download-path", "d", "", "Path new release should be downloaded to (required)")
    rootCmd.Flags().StringVarP(&extractPath, "extract-path", "e", "", "Path release should be extracted to (required)")
    rootCmd.MarkFlagRequired("download-path")
    rootCmd.MarkFlagRequired("extract-path")

    if err := rootCmd.Execute(); err != nil {
        system.Exit(1)
    }
}
