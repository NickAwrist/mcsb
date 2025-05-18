package internal

import (
	"bytes"
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

	fmt.Printf("\nCreating server %s with framework %s and version %s\n", serverName, framework, version.ID)

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
	err = runServer(filename)
	if err != nil {
		fmt.Printf("Failed to generate initial configuration: %v\n", err)
		os.Exit(1)
	}
	bar.Add(1)
	time.Sleep(500 * time.Millisecond)

	// Step 4: Accept EULA
	bar.Describe("Step 4/5: Accepting EULA")
	time.Sleep(2 * time.Second)

	eulaErr := acceptEula()
	if eulaErr != nil {
		fmt.Printf("Failed to accept EULA: %v\n", eulaErr)
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

func runServer(filename string) error {
	cmd := exec.Command("java", "-jar", filename, "nogui")

	// Capture stderr to check for Java version errors
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return err
	}

	errChan := make(chan error, 1)
	fileChan := make(chan bool, 1)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			if _, err := os.Stat("eula.txt"); err == nil {
				fileChan <- true
				return
			}
			if strings.Contains(stderr.String(), "Unsupported Java") {
				errChan <- fmt.Errorf("%s", stderr.String())
				return
			}
			// Check if process exited with error
			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				errChan <- fmt.Errorf("server process exited unexpectedly: %s", stderr.String())
				return
			}
		}
	}()

	select {
	case err := <-errChan:
		fmt.Printf("Error during server initialization: %v\n", err)
		return err
	case <-fileChan:
		// Successfully generated files
	case <-time.After(30 * time.Second):
		fmt.Println("Timeout waiting for server initialization")
		return fmt.Errorf("timeout waiting for server initialization")
	}
	cmd.Process.Kill()
	return nil
}

func acceptEula() error {
	// Edit the eula.txt file
	eulaFile, err := os.OpenFile("eula.txt", os.O_RDWR, 0644)
	if err != nil {
		eulaFile, err = os.Create("eula.txt")
		if err != nil {
			return fmt.Errorf("error creating eula.txt: %w", err)
		}
		_, err = eulaFile.WriteString("eula=true")
		if err != nil {
			return fmt.Errorf("error writing to eula.txt: %w", err)
		}
		return nil
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
