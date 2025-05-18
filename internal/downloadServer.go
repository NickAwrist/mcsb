package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// VersionMetadata represents the structure of the version metadata JSON
type VersionMetadata struct {
	Downloads struct {
		Server struct {
			URL string `json:"url"`
		} `json:"server"`
	} `json:"downloads"`
}

func DownloadVanillaServer(version Version, downloadDir string) string {
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

	filename := version.ID
	for i, char := range filename {
		if char == '.' {
			filename = filename[:i] + "-" + filename[i+1:]
		}
	}
	filename = filename + ".jar"

	DownloadFileURL(metadata.Downloads.Server.URL, filename, downloadDir)

	return filename
}

func DownloadPaperServer(version Version, downloadDir string) string {
	resp, err := http.Get("https://api.papermc.io/v2/projects/paper/versions/" + version.ID + "/builds")
	if err != nil {
		fmt.Println("Error fetching version metadata:", err)
		os.Exit(1)
	}

	var decodedBuilds PaperBuildOutput
	err = json.NewDecoder(resp.Body).Decode(&decodedBuilds)
	if err != nil {
		fmt.Println("Error decoding version metadata:", err)
		os.Exit(1)
	}

	if len(decodedBuilds.Builds) == 0 {
		fmt.Printf("\nNo builds found for version %s\n", version.ID)
		os.Exit(1)
	}

	latestBuild := decodedBuilds.Builds[len(decodedBuilds.Builds)-1]

	url := fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s/builds/%d/downloads/%s", version.ID, latestBuild.Build, latestBuild.Downloads.Application.Name)

	filename := latestBuild.Downloads.Application.Name
	strings.ReplaceAll(filename, ".", "-")
	strings.Replace(filename, "-jar", ".jar", 1)

	DownloadFileURL(url, filename, downloadDir)

	return filename
}
