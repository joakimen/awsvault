package cmd

import (
	"fmt"
	"github.com/joakimen/awsvault/pkg/aws"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles in ~/.aws/config",
	Run: func(cmd *cobra.Command, args []string) {
		listFn()
	},
}

func listFn() {
	conf, err := aws.ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't read configuration file: %v\n", err)
		os.Exit(1)
	}

	profiles := aws.ParseAWSProfiles(conf)
	if profiles == nil {
		fmt.Fprintf(os.Stderr, "no assumable profiles found in configuration file\n")
		os.Exit(1)
	}

	for _, profile := range profiles {
		fmt.Println(profile.Name)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
