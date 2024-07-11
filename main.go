package main

import (
	"fmt"
	
	"updater/internal/utils/error"
	"updater/internal/commands/update"
    "updater/internal/utils/filedeleter"
	"updater/internal/utils/unzipper"
	"updater/internal/utils/system"
	"updater/internal/utils/versionchecker"
	
	"github.com/inancgumus/screen"
	"github.com/spf13/cobra"
)

type Options struct {
    RepoUrl string
    CheckVersion bool
    TempDownloadPath string
    ExtractionPath string
    HasNewVersion bool
    ReleaseDownloadUrl string
    ReleaseName string
    ReleaseVersion string
    InstalledVersionPath string
    DownloadedFile string
    UpdateStatus string
}

func main() {
	screen.Clear()
	screen.MoveTopLeft()
	
	var opts Options

	rootCmd := &cobra.Command{
		Use: "vue-updater",
		Short: "A CLI tool to update vue ui torrent frontend",
		Run: func(cmd *cobra.Command, args []string) {
			checkVersionAndUpdate(&opts)
			
			fmt.Println()
			fmt.Println(opts.UpdateStatus)
			fmt.Println()
		},
	}

	// Define flags with clearer descriptions and ensure required fields are marked as such.
	flags := rootCmd.Flags()
	flags.BoolVarP(&opts.CheckVersion, "check-version", "c", true, "Check to see if the versions are different before downloading")
	flags.StringVarP(&opts.TempDownloadPath, "temp-download-path", "t", "", "Temporary path to download the new file (required)")
	flags.StringVarP(&opts.ExtractionPath, "extraction-path", "e", "", "Path to extract the release file to (required)")
	flags.StringVarP(&opts.RepoUrl, "repo-url", "r", "https://api.github.com/repos/VueTorrent/VueTorrent/releases/latest", "Repo Release Url you wish to download from (optional)")
	rootCmd.MarkFlagRequired("temp-download-path")
	rootCmd.MarkFlagRequired("extraction-path")

	if err := rootCmd.Execute(); err != nil {
		system.Exit(1)
	}
}

func checkVersionAndUpdate(opts *Options) {
    if opts.CheckVersion {
        versionPath := fmt.Sprintf("%s/%s", opts.ExtractionPath, "vuetorrent/version.txt")
        isNewVersion, name, url, ver, err := versionchecker.IsNewVersion(opts.RepoUrl, versionPath)
        if err != nil {
            error.Log(err, "Failed to check for new version")
            return
        }
        
        opts.HasNewVersion = isNewVersion
        opts.ReleaseName = name
        opts.ReleaseDownloadUrl = url
        opts.ReleaseVersion = ver
    } else {
        opts.UpdateStatus = "Already on latest version"
        return // Added to avoid unnecessary execution of subsequent code when not needed.
    }
    
    if opts.HasNewVersion {
        file, err := update.GetUpdate(opts.ReleaseName, opts.ReleaseDownloadUrl, opts.TempDownloadPath)
        if err != nil {
            error.Log(err, "Failed to update version")
            return
        }
        
        opts.DownloadedFile = file
        if err := unzipper.UnzipWithProgress(opts.DownloadedFile, opts.ExtractionPath); err != nil {
            error.Log(err, "Failed to extract file")
            return
        }
        
        if err := filedeleter.DeleteFileWithProgress(opts.DownloadedFile); err != nil {
            error.Log(err, "Failed to delete file")
        } else {
            opts.UpdateStatus = fmt.Sprintf("Updated successfully to Version: %s", opts.ReleaseVersion)
        }
    } else {
        opts.UpdateStatus = "Already on latest version"
    }
}