package cmd

import (
	"fmt"
	"github.com/joakimen/awsvault/pkg/aws"
	"github.com/joakimen/awsvault/pkg/fs"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	chromePath  string = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	firefoxPath string = "/Applications/Firefox.app/Contents/MacOS/firefox"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open a browser window and login using temporary credentials",
	Run: func(cmd *cobra.Command, args []string) {
		loginFn(cmd)
	},
}

func loginFn(cmd *cobra.Command) {

	browser, _ := cmd.Flags().GetString("browser")
	if browser != "chrome" && browser != "firefox" {
		fmt.Fprintf(os.Stderr, "Invalid browser: %s\n", browser)
		os.Exit(1)
	}

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
	output, err := command.CombinedOutput()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting login link using aws-vault login: %v\n%s\n", err, output)
		os.Exit(1)
	}

	loginLink := string(output)

	if browser == "chrome" {
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
			fmt.Fprintf(os.Stderr, "error starting browser '%s': %v\n", browser, err)
			os.Exit(1)
		}

	} else {
		// need this before passing link to tab container extension
		urlEncodedLoginLink := url.QueryEscape(loginLink)

		var groupColor string
		if strings.Contains(profile.Name, "-prod") {
			groupColor = "red"
		} else if strings.Contains(profile.Name, "-staging") {
			groupColor = "yellow"
		} else {
			groupColor = "green"
		}

		containerTabCommand := fmt.Sprintf("ext+container:name=%s&color=%s&url=%s", profile.Name, groupColor, urlEncodedLoginLink)
		firefox := exec.Command(firefoxPath, containerTabCommand)
		err = firefox.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error starting browser '%s': %v\n", browser, err)
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("browser", "b", "firefox", "Browser to use (chrome, firefox, default)")
}
