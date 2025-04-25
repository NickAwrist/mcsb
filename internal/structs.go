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

type manifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []Version `json:"versions"`
}
