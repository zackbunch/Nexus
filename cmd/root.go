package cmd

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Default values - you can replace these with your actual defaults
	defaultNexusURL = "https://your-nexus-instance.com"
	defaultUsername = "service-account-username"
	defaultToken    = "your-default-token"

	// CLI flags
	sourceDir      string
	nexusRepo      string
	nexusURL       string
	nexusUsername  string
	nexusToken     string
	nexusDirectory string
)

var rootCmd = &cobra.Command{
	Use:   "nexus-uploader",
	Short: "Upload files to Nexus repository",
	Long:  `A CLI tool to upload files from a local directory to Nexus repository.`,
	RunE:  upload,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&sourceDir, "source", "s", "", "Source directory containing files to upload (required)")
	rootCmd.Flags().StringVarP(&nexusRepo, "repository", "r", "", "Nexus repository name (required)")
	rootCmd.Flags().StringVarP(&nexusURL, "url", "u", defaultNexusURL, "Nexus server URL")
	rootCmd.Flags().StringVar(&nexusUsername, "username", defaultUsername, "Nexus username")
	rootCmd.Flags().StringVar(&nexusToken, "token", defaultToken, "Nexus token")
	rootCmd.Flags().StringVarP(&nexusDirectory, "directory", "d", "", "Target directory in Nexus (required)")

	rootCmd.MarkFlagRequired("source")
	rootCmd.MarkFlagRequired("repository")
	rootCmd.MarkFlagRequired("directory")
}

func upload(cmd *cobra.Command, args []string) error {
	// Validate source directory
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", sourceDir)
	}

	// Walk through the source directory
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Upload the file
		return uploadFile(path, relPath)
	})
}

func uploadFile(localPath, relPath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Determine content type
	contentType := mime.TypeByExtension(filepath.Ext(localPath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Construct the upload URL
	uploadPath := strings.TrimPrefix(filepath.Join(nexusDirectory, relPath), "/")
	uploadURL := fmt.Sprintf("%s/repository/%s/%s",
		strings.TrimSuffix(nexusURL, "/"),
		nexusRepo,
		uploadPath)

	// Create request
	req, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	req.SetBasicAuth(nexusUsername, nexusToken)

	// Perform the upload
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s - %s", resp.Status, string(body))
	}

	fmt.Printf("Successfully uploaded: %s\n", relPath)
	return nil
}
