package cmd

import (
	"fmt"
	"github.com/joakimen/awsvault/pkg/aws"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Start a subshell with temporary credentials",
	Run: func(cmd *cobra.Command, args []string) {
		execFn()
	},
}

func execFn() {
	configFile, err := aws.ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't read configuration file: %v\n", err)
		os.Exit(1)
	}

	profiles := aws.ParseAWSProfiles(configFile)
	if profiles == nil {
		fmt.Fprintf(os.Stderr, "no assumable profiles found in configuration file\n")
		os.Exit(1)
	}

	selectedProfile := aws.FuzzySelectAWSProfile(profiles)

	command := exec.Command("aws-vault", "exec", selectedProfile.Name)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err = command.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running aws-vault exec: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(execCmd)
}
