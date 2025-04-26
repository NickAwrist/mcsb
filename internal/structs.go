package internal

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type ColoredOption struct {
	OptionText string
	Color      string
}

var Templates = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   fmt.Sprintf("%s {{ .Color | cyan }}", promptui.IconSelect),
	Inactive: " {{ .Color }}",
	Selected: fmt.Sprintf("%s {{ .Color | green }}", promptui.IconGood),
}

type Version struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type MojangManifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []Version `json:"versions"`
}

type PaperManifest struct {
	ProjectName string   `json:"project_name"`
	Versions    []string `json:"versions"`
}

type PaperBuildOutput struct {
	Version string       `json:"version"`
	Builds  []PaperBuild `json:"builds"`
}

type PaperBuild struct {
	Build     int    `json:"build"`
	Channel   string `json:"channel"`
	Downloads struct {
		Application struct {
			Name string `json:"name"`
		} `json:"application"`
	} `json:"downloads"`
}
