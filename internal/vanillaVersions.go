package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetVanillaVersions() []Version {

	resp, err := http.Get("https://launchermeta.mojang.com/mc/game/version_manifest.json")
	if err != nil {
		fmt.Println("Error fetching versions:", err)
		return nil
	}
	defer resp.Body.Close()

	var responseManifest manifest
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
