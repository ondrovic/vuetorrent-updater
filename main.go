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

type Options struct {
    Url string
    DownloadPath string
    ExtractionPath string
    DownloadedFile string
}

func main() {
    screen.Clear()
	screen.MoveTopLeft()
    
    var opts Options

    rootCmd := &cobra.Command {
        Use: "vue-updater",
        Short: "A CLI tool to update vue ui torrent frontend",
        Run: func(cmd *cobra.Command, args []string) {
            downloadedFile, err := update.RunCommand(opts.DownloadPath, opts.Url)
            if err != nil {
                error.Log(err, "Failed to download file")
                return
            }
            opts.DownloadedFile = downloadedFile

            if err := unzipper.UnzipWithProgress(opts.DownloadedFile, opts.ExtractionPath); err != nil {
                error.Log(err, "Failed to unzip archive")
                return
            }

            if err := filedeleter.DeleteFileWithProgress(opts.DownloadedFile); err != nil {
                error.Log(err, "Failed to remove downloaded archive")
                return
            }

            fmt.Println("Update completed successfully")
        },
    }

    rootCmd.Flags().StringVarP(&opts.Url, "url", "u", "", "Url you wish to download from")
    rootCmd.Flags().StringVarP(&opts.DownloadPath, "download-path", "d", "", "Path new release should be downloaded to (required)")
    rootCmd.Flags().StringVarP(&opts.ExtractionPath, "extract-path", "e", "", "Path release should be extracted to (required)")
    rootCmd.MarkFlagRequired("download-path")
    rootCmd.MarkFlagRequired("extract-path")
    rootCmd.MarkFlagRequired("url")

    if err := rootCmd.Execute(); err != nil {
        system.Exit(1)
    }
}
