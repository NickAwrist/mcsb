package cmd

import (
	"fmt"
	"mcsb-cli/internal"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Minecraft Server",
	Long:  "Start the creation process of a new Minecraft Server",
	Run: func(cmd *cobra.Command, args []string) {
		framework := getFramework()
		fmt.Printf("Selected framework: %s\n", framework)
		version := getServerVersion(framework)
		fmt.Printf("Selected version: %s\n", version)

		serverName := getServerName()
		fmt.Printf("Selected server name: %s\n", serverName)

		serverPort := getServerPort()
		fmt.Printf("Selected server port: %d\n", serverPort)

		eula := acceptEula()
		fmt.Printf("EULA accepted: %t\n", eula)

		if !eula {
			fmt.Println("EULA not accepted, exiting...")
			os.Exit(1)
		}

		internal.CreateServer(framework, version, serverName, serverPort)

		fmt.Println("Server created successfully!")

	},
}

func getFramework() string {
	coloredOptions := []internal.ColoredOption{
		{OptionText: "PaperMC", Color: color.CyanString("PaperMC (recommended)")},
		{OptionText: "Vanilla", Color: color.GreenString("Vanilla")},
		{OptionText: "Spigot", Color: color.YellowString("Spigot")},
	}

	frameworkPrompt := promptui.Select{
		Label:     "Select the server framework",
		Items:     coloredOptions,
		Templates: internal.Templates,
	}

	frameworkResult, _, err := frameworkPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return coloredOptions[frameworkResult].OptionText
}

func getServerVersion(framework string) internal.Version {
	if framework == "Vanilla" {
		versions := internal.GetVanillaVersions()

		coloredVersions := []internal.ColoredOption{}
		for _, v := range versions {
			coloredVersions = append(coloredVersions, internal.ColoredOption{OptionText: v.ID, Color: color.YellowString(v.ID)})
		}

		latestVersion := coloredVersions[0].OptionText + " (latest)"
		coloredVersions[0].Color = color.GreenString(latestVersion)
		coloredVersions[0].OptionText = latestVersion

		versionPrompt := promptui.Select{
			Label:     "Select the server version",
			Items:     coloredVersions,
			Templates: internal.Templates,
		}

		versionResult, _, err := versionPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}

		return versions[versionResult]
	}

	return internal.Version{}
}

func getServerName() string {
	serverNamePrompt := promptui.Prompt{
		Label: "Enter your server name",
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
		Label: "Enter the server port",
		Validate: func(input string) error {
			port, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid port: %v", err)
			}
			if port < 1 || port > 65535 {
				return fmt.Errorf("port must be between 1 and 65535")
			}
			return nil
		},
	}

	serverPortResult, err := serverPortPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	serverPort, err := strconv.Atoi(serverPortResult)
	if err != nil {
		fmt.Printf("Invalid port: %v\n", err)
		os.Exit(1)
	}

	return serverPort
}

func acceptEula() bool {

	coloredOptions := []internal.ColoredOption{
		{OptionText: "Yes", Color: color.GreenString("Yes")},
		{OptionText: "No", Color: color.RedString("No")},
	}

	eulaPrompt := promptui.Select{
		Label:     "Accept the EULA",
		Items:     coloredOptions,
		Templates: internal.Templates,
	}

	eulaResult, _, err := eulaPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return eulaResult == 0
}

func Execute() {
	if err := createCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
