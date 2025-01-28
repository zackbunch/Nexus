package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	nexusURL   string
	username   string
	token      string
	repository string
	sourceDir  string
	targetDir  string
	pathPrefix string
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "nexus-cli",
	Short: "A CLI tool for interacting with Nexus Repository",
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files to Nexus repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement upload logic
		return nil
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download files from Nexus repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := NewClient(nexusURL, username, token)
		if err != nil {
			return fmt.Errorf("failed to create nexus client: %w", err)
		}

		if configFile != "" {
			config, err := LoadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			for _, repo := range config.Repositories {
				for _, dir := range repo.Directories {
					err := client.RepositoryService.DownloadRepositoryAssets(
						repo.Name,
						repo.LocalPath,
						dir,
					)
					if err != nil {
						return fmt.Errorf("failed to download from %s/%s: %w", repo.Name, dir, err)
					}
				}
			}
			return nil
		}

		return client.RepositoryService.DownloadRepositoryAssets(
			repository,
			targetDir,
			pathPrefix,
		)
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&nexusURL, "url", "", "Nexus server URL")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Nexus username")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Nexus token")
	rootCmd.PersistentFlags().StringVar(&repository, "repository", "", "Repository name")

	// Upload command flags
	uploadCmd.Flags().StringVar(&sourceDir, "source", "", "Source directory to upload")
	uploadCmd.Flags().StringVar(&targetDir, "target", "", "Target directory in Nexus")

	// Download command flags
	downloadCmd.Flags().StringVar(&targetDir, "target", "", "Local directory to download to")
	downloadCmd.Flags().StringVar(&pathPrefix, "prefix", "", "Path prefix to filter assets")
	downloadCmd.Flags().StringVar(&configFile, "config", "", "Path to TOML config file")

	rootCmd.AddCommand(uploadCmd, downloadCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
