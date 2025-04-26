package internal

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
)

func DownloadFileURL(url string, filename string, downloadDir string) {

	println("Downloading " + url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	serverResp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading server:", err)
		os.Exit(1)
	}
	defer serverResp.Body.Close()

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
		fmt.Println("Error saving .jar file:", err)
		os.Exit(1)
	}
}
