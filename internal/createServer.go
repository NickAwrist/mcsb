package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func CreateServer(framework string, version Version, serverName string, serverPort int) {
	// Setup progress tracking
	steps := []string{
		"Creating server directory",
		"Downloading server files",
		"Generating initial configuration",
		"Accepting EULA",
		"Customizing server properties",
	}

	totalSteps := len(steps)
	bar := progressbar.NewOptions(totalSteps,
		progressbar.OptionSetDescription("Creating server"),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// Step 1: Create directory
	bar.Describe("Step 1/5: Creating server directory")
	homeDir, err := os.UserHomeDir()
	if err != nil {
	}
	desktopDir := filepath.Join(homeDir, "Desktop")
	downloadDir := filepath.Join(desktopDir, serverName)
	os.MkdirAll(downloadDir, 0755)
	os.Chdir(downloadDir)
	bar.Add(1)
	time.Sleep(500 * time.Millisecond)

	// Step 2: Download server
	bar.Describe("Step 2/5: Downloading server files")
	var filename string
	switch framework {
	case "Vanilla":
		filename = DownloadVanillaServer(version, downloadDir)
	case "PaperMC":
		filename = DownloadPaperServer(version, downloadDir)
	}
	bar.Add(1)
	time.Sleep(500 * time.Millisecond)

	// Step 3: Generate initial configuration
	bar.Describe("Step 3/5: Generating initial configuration")
	cmd := exec.Command("java", "-jar", filename, "nogui")
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}

	// Wait for server to generate files
	time.Sleep(8 * time.Second)
	cmd.Process.Kill()
	bar.Add(1)
	time.Sleep(500 * time.Millisecond)

	// Step 4: Accept EULA
	bar.Describe("Step 4/5: Accepting EULA")
	time.Sleep(2 * time.Second)

	// Try multiple times to access the file
	var eulaErr error
	for attempts := 0; attempts < 3; attempts++ {
		eulaErr = acceptEula()
		if eulaErr == nil {
			break
		}
		fmt.Printf("Retrying EULA acceptance (%d/3)...\n", attempts+1)
		time.Sleep(2 * time.Second)
	}

	if eulaErr != nil {
		fmt.Printf("Failed to accept EULA after multiple attempts: %v\n", eulaErr)
		os.Exit(1)
	}

	bar.Add(1)
	time.Sleep(500 * time.Millisecond)

	// Step 5: Edit server properties
	bar.Describe("Step 5/5: Customizing server properties")
	time.Sleep(2 * time.Second)

	// Try multiple times to access the properties file
	var propsErr error
	for attempts := 0; attempts < 3; attempts++ {
		propsErr = editServerProperties(serverPort, serverName)
		if propsErr == nil {
			break
		}
		fmt.Printf("Retrying properties edit (%d/3)...\n", attempts+1)
		time.Sleep(2 * time.Second)
	}

	if propsErr != nil {
		fmt.Printf("Failed to edit properties after multiple attempts: %v\n", propsErr)
		os.Exit(1)
	}

	bar.Add(1)
	time.Sleep(500 * time.Millisecond)

}

func acceptEula() error {
	// Edit the eula.txt file
	eulaFile, err := os.OpenFile("eula.txt", os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening eula.txt: %w", err)
	}
	defer eulaFile.Close()

	// Replace the eula=false with eula=true
	content, err := io.ReadAll(eulaFile)
	if err != nil {
		return fmt.Errorf("error reading eula.txt: %w", err)
	}
	contentStr := strings.Replace(string(content), "eula=false", "eula=true", 1)
	_, err = eulaFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error seeking in eula.txt: %w", err)
	}
	err = eulaFile.Truncate(0)
	if err != nil {
		return fmt.Errorf("error truncating eula.txt: %w", err)
	}
	_, err = eulaFile.WriteString(contentStr)
	if err != nil {
		return fmt.Errorf("error writing to eula.txt: %w", err)
	}

	return nil
}

func editServerProperties(serverPort int, serverName string) error {
	// Edit the server.properties file
	serverPropertiesFile, err := os.OpenFile("server.properties", os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening server.properties: %w", err)
	}
	defer serverPropertiesFile.Close()

	// Replace the server-port and server-name with the given values
	content, err := io.ReadAll(serverPropertiesFile)
	if err != nil {
		return fmt.Errorf("error reading server.properties: %w", err)
	}
	contentStr := strings.Replace(string(content), "server-port=25565", fmt.Sprintf("server-port=%d", serverPort), 1)
	contentStr = strings.Replace(contentStr, "motd=A Minecraft Server", fmt.Sprintf("motd=%s", serverName), 1)
	_, err = serverPropertiesFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error seeking in server.properties: %w", err)
	}
	err = serverPropertiesFile.Truncate(0)
	if err != nil {
		return fmt.Errorf("error truncating server.properties: %w", err)
	}
	_, err = serverPropertiesFile.WriteString(contentStr)
	if err != nil {
		return fmt.Errorf("error writing to server.properties: %w", err)
	}

	return nil
}
