package aws

import (
	"encoding/json"
	"fmt"
	fz "github.com/ktr0731/go-fuzzyfinder"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
	Name           string
	Region         string
	RoleArn        string
	MfaSerial      string
	SourceProfile  string
	IncludeProfile string
	Output         string
}

func ReadConfig() (*ini.File, error) {
	configfilePath := filepath.Join(os.Getenv("HOME"), ".aws", "config")
	conf, err := ini.Load(configfilePath)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func ParseAWSProfiles(awsConfigFile *ini.File) []Profile {
	var awsProfiles []Profile
	profilePrefix := "profile "
	for _, section := range awsConfigFile.Sections() {
		if strings.HasPrefix(section.Name(), profilePrefix) {
			profile := Profile{Name: strings.Replace(section.Name(), profilePrefix, "", 1)}
			for _, key := range section.Keys() {
				switch key.Name() {
				case "region":
					profile.Region = key.Value()
				case "role_arn":
					profile.RoleArn = key.Value()
				case "mfa_serial":
					profile.MfaSerial = key.Value()
				case "source_profile":
					profile.SourceProfile = key.Value()
				case "include_profile":
					profile.IncludeProfile = key.Value()
				case "output":
					profile.Output = key.Value()
				}
			}

			// if a profile is missing RoleArn, it's a meta-profile and should be ignored
			if profile.RoleArn != "" {
				awsProfiles = append(awsProfiles, profile)
			}
		}
	}
	if len(awsProfiles) > 0 {
		return awsProfiles
	}
	return nil
}

func FuzzySelectAWSProfile(existingFiles []Profile) Profile {

	renderFunc := func(selectedIndex int) string {
		return existingFiles[selectedIndex].Name
	}

	previewFunc := func(selectedIndex, width, height int) string {
		if selectedIndex == -1 {
			return ""
		}

		s, _ := json.MarshalIndent(existingFiles[selectedIndex], "", "  ")
		return string(s)
	}

	idx, err := fz.Find(existingFiles, renderFunc, fz.WithPreviewWindow(previewFunc))
	if err != nil {
		fmt.Println("no profile selected")
		os.Exit(0)
	}

	return existingFiles[idx]
}
