package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

// VersionMetadata represents the structure of the version metadata JSON
type VersionMetadata struct {
	Downloads struct {
		Server struct {
			URL string `json:"url"`
		} `json:"server"`
	} `json:"downloads"`
}

func DownloadVanillaServer(version Version) string {
	fmt.Printf("\nDownloading Minecraft server version %s...\n", version.ID)

	resp, err := http.Get(version.URL)
	if err != nil {
		fmt.Println("Error fetching version metadata:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var metadata VersionMetadata
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		fmt.Println("Error decoding version metadata:", err)
		os.Exit(1)
	}

	serverResp, err := http.Get(metadata.Downloads.Server.URL)
	if err != nil {
		fmt.Println("Error downloading server:", err)
		os.Exit(1)
	}
	defer serverResp.Body.Close()

	filename := version.ID
	for i, char := range filename {
		if char == '.' {
			filename = filename[:i] + "-" + filename[i+1:]
		}
	}
	filename = filename + ".jar"

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating server.jar file:", err)
		os.Exit(1)
	}
	defer f.Close()

	bar := progressbar.NewOptions64(
		serverResp.ContentLength,
		progressbar.OptionSetDescription("Downloading JAR"),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	_, err = io.Copy(io.MultiWriter(f, bar), serverResp.Body)
	if err != nil {
		fmt.Println("Error saving server.jar file:", err)
		os.Exit(1)
	}

	fmt.Printf("Minecraft server version %s downloaded successfully!\n", version.ID)
	return filename
}
