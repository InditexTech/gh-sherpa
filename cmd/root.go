package cmd

import (
	"os"
	"strings"

	"github.com/InditexTech/gh-sherpa/cmd/create_branch"
	"github.com/InditexTech/gh-sherpa/cmd/create_pull_request"
	"github.com/InditexTech/gh-sherpa/internal/config"
	"github.com/InditexTech/gh-sherpa/internal/logging"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "sherpa",
	Short:         "Interact with Inditex Sherpa",
	Long:          "GitHub CLI Sherpa Extension - Interact with the Inditex Sherpa service from the command line.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Ignore initialization when the command is help
		if strings.HasPrefix(cmd.Use, "help") {
			return nil
		}

		isInteractive := !useDefaultValues
		return config.Initialize(isInteractive)
	},
}

const versionTemplate string = `GitHub CLI Sherpa Extension - Version: {{ printf "%s\n" .Version }}`

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		logging.Error(err.Error())
		os.Exit(1)
	}
}

var useDefaultValues bool

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetVersionTemplate(versionTemplate)

	rootCmd.PersistentFlags().BoolVarP(&useDefaultValues, "yes", "y", false, "use the default proposed fields")

	rootCmd.AddCommand(create_branch.Command)
	rootCmd.AddCommand(create_pull_request.Command)
}

func SetVersion(version string) {
	rootCmd.Version = version
}
