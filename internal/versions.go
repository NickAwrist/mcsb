package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

func GetVanillaVersions() []Version {

	resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		fmt.Println("Error fetching versions:", err)
		return nil
	}
	defer resp.Body.Close()

	var responseManifest MojangManifest
	err = json.NewDecoder(resp.Body).Decode(&responseManifest)
	if err != nil {
		fmt.Println("Error decoding versions:", err)
		return nil
	}

	var versions []Version
	for _, v := range responseManifest.Versions {
		if v.Type == "release" {
			versions = append(versions, v)
		}
	}

	return versions
}

func GetPaperVersions() []Version {

	resp, err := http.Get("https://api.papermc.io/v2/projects/paper")
	if err != nil {
		fmt.Println("Error fetching versions:", err)
		return nil
	}
	defer resp.Body.Close()

	var responseManifest PaperManifest
	err = json.NewDecoder(resp.Body).Decode(&responseManifest)
	if err != nil {
		fmt.Println("Error decoding versions:", err)
		return nil
	}

	slices.Reverse(responseManifest.Versions)
	var versions []Version
	for i, v := range responseManifest.Versions {
		versionType := "default"
		if i == 0 {
			resp, err := http.Get("https://api.papermc.io/v2/projects/paper/versions/" + v + "/builds")
			if err != nil {
			}

			var decodedBuilds PaperBuildOutput
			err = json.NewDecoder(resp.Body).Decode(&decodedBuilds)
			if err != nil {
			}

			latestBuild := decodedBuilds.Builds[len(decodedBuilds.Builds)-1]

			versionType = latestBuild.Channel
		}

		latestVersion := Version{
			ID:   v,
			Type: versionType,
			URL:  "TBD",
		}

		versions = append(versions, latestVersion)
	}

	return versions
}
