package cmd

import (
	"errors"
	"fmt"
	"log"
	"mcsb-cli/internal"
	"mcsb-cli/util"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Minecraft Server",
	Long:  "Start the creation process of a new Minecraft Server",
	Run: func(cmd *cobra.Command, args []string) {
		// Get which framework the user wants to use (Vanilla, Paper)
		framework, _ := cmd.Flags().GetString("framework")
		if framework != "" {
			if strings.ToLower(framework) == "paper" || strings.ToLower(framework) == "papermc" {
				framework = "PaperMC"
			} else if strings.ToLower(framework) == "vanilla" {
				framework = "Vanilla"
			} else {
				fmt.Printf("Invalid framework '%s'. Using interactive prompt instead.\n", framework)
				framework = ""
			}
		}
		if framework == "" {
			framework = getFramework()
		}

		// What version of Minecraft the user wants
		version, _ := cmd.Flags().GetString("version")
		if version == "" {
			version = getServerVersion(framework)
		}

		// The name of the server
		serverName, _ := cmd.Flags().GetString("name")
		if serverName == "" {
			serverName = getServerName()
		}

		// Desired port
		serverPort, _ := cmd.Flags().GetInt("port")
		if serverPort == 0 {
			serverPort = getServerPort()
		}

		// Whether they want to accept the EULA
		acceptEula, _ := cmd.Flags().GetBool("eula")
		if !acceptEula {
			acceptEulaPrompt()
		}

		// With all the information, create the server directory and download the files
		internal.CreateServer(framework, internal.Version{ID: version}, serverName, serverPort)

		fmt.Println("\nServer created successfully!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("framework", "f", "", "Server framework (PaperMC or Vanilla)")
	createCmd.Flags().StringP("version", "v", "", "Minecraft version")
	createCmd.Flags().StringP("name", "n", "", "Server name")
	createCmd.Flags().IntP("port", "p", 25565, "Server port")
	createCmd.Flags().BoolP("eula", "e", false, "Accept EULA")
}

func getFramework() string {
	coloredOptions := []internal.ColoredOption{
		{OptionText: "PaperMC", Color: color.CyanString("PaperMC (recommended)")},
		{OptionText: "Vanilla", Color: color.GreenString("Vanilla")},
	}

	frameworkPrompt := promptui.Select{
		Label:     "Select the server framework",
		Items:     coloredOptions,
		Templates: internal.Templates,
		Stdout:    util.NoBellStdout,
	}

	frameworkResult, _, err := frameworkPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return coloredOptions[frameworkResult].OptionText
}

func getServerVersion(framework string) string {
	var versions []internal.Version
	switch framework {
	case "Vanilla":
		versions = internal.GetVanillaVersions()
	case "PaperMC":
		versions = internal.GetPaperVersions()
	}

	if len(versions) == 0 {
		log.Fatalf("no versions found for %s", framework)
	}

	var coloredVersions []internal.ColoredOption
	for _, v := range versions {
		coloredVersions = append(coloredVersions, internal.ColoredOption{OptionText: v.ID, Color: color.YellowString(v.ID)})
	}

	latestVersionIndex := 0
	if framework == "PaperMC" && versions[0].Type == "experimental" {
		experimentalVersion := coloredVersions[0].OptionText + "(experimental)"
		coloredVersions[0].Color = color.RedString(experimentalVersion)
		coloredVersions[0].OptionText = experimentalVersion
		latestVersionIndex = 1
	}

	latestVersion := coloredVersions[latestVersionIndex].OptionText + " (latest)"
	coloredVersions[latestVersionIndex].Color = color.GreenString(latestVersion)
	coloredVersions[latestVersionIndex].OptionText = latestVersion

	versionPrompt := promptui.Select{
		Label:     "Select the server version",
		Items:     coloredVersions,
		Templates: internal.Templates,
		Stdout:    util.NoBellStdout,
	}

	versionResult, _, err := versionPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return versions[versionResult].ID
}

func getServerName() string {
	serverNamePrompt := promptui.Prompt{
		Label:       "Enter your server name",
		HideEntered: true,
	}

	serverNameResult, err := serverNamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return serverNameResult
}

func getServerPort() int {
	serverPortPrompt := promptui.Prompt{
		Label: "Enter the server port (default: 25565)",
		Validate: func(input string) error {
			if input == "" {
				return nil
			}
			port, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid port: %v", err)
			}
			if port < 1 || port > 65535 {
				return fmt.Errorf("port must be between 1 and 65535")
			}
			return nil
		},
		HideEntered: true,
	}

	serverPortResult, err := serverPortPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	if serverPortResult == "" {
		serverPortResult = "25565"
	}

	serverPort, err := strconv.Atoi(serverPortResult)
	if err != nil {
		fmt.Printf("Invalid port: %v\n", err)
		os.Exit(1)
	}

	return serverPort
}

func acceptEulaPrompt() {
	eulaPrompt := promptui.Prompt{
		Label:       "Do you wish to accept Mojang's EULA?",
		IsConfirm:   true,
		Default:     "Y",
		HideEntered: true,
	}
	_, err := eulaPrompt.Run()

	switch {
	case errors.Is(err, promptui.ErrAbort):
		fmt.Println("You did not accept Mojang's EULA, exiting server creation.")
		os.Exit(0)
	case errors.Is(err, promptui.ErrInterrupt):
		fmt.Println("Exiting server creation.")
	}
}

func ExecuteCreate() {
	if err := createCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
