package cmd

import (
	"fmt"
	"github.com/joakimen/awsvault/pkg/aws"
	"github.com/spf13/cobra"
	"os"
)

// selectCmd represents the select command
var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select a single AWS profile interactively",
	Run: func(cmd *cobra.Command, args []string) {
		selectFn()
	},
}

func selectFn() {
	config, err := aws.ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't read configuration file: %v\n", err)
		os.Exit(1)
	}

	profiles := aws.ParseAWSProfiles(config)
	if profiles == nil {
		fmt.Fprintf(os.Stderr, "no assumable profiles found in configuration file\n")
		os.Exit(1)
	}

	profile := aws.FuzzySelectAWSProfile(profiles)
	fmt.Println(profile.Name)
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
