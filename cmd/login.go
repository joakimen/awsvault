package cmd

import (
	"fmt"
	"github.com/joakimen/awsvault/pkg/aws"
	"github.com/joakimen/awsvault/pkg/fs"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open a browser window and login using temporary credentials",
	Run: func(cmd *cobra.Command, args []string) {
		loginFn()
	},
}

func loginFn() {

	var (
		chromePath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	)

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

	profile := aws.FuzzySelectAWSProfile(profiles)

	command := exec.Command("aws-vault", "login", profile.Name, "--stdout", "--prompt", "osascript")
	stdout, err := command.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting login link using aws-vault login: %v\n", err)
		os.Exit(1)
	}

	loginLink := string(stdout)
	xdgDataDir := filepath.Join(fs.XDGDataDir(), "aws_chrome", profile.Name)
	xdgCacheDir := filepath.Join(fs.XDGCacheDir(), "aws_chrome", profile.Name)

	chrome := exec.Command(chromePath,
		"--no-first-run",
		"--start-maximized",
		"--user-data-dir="+xdgDataDir,
		"--disk-cache-dir="+xdgCacheDir,
		loginLink)

	err = chrome.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error starting browser: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
