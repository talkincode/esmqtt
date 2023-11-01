package assets

import (
	_ "embed"
	"regexp"
)

//go:embed buildinfo.txt
var BuildInfo string

var defaultBuildVer = "Latest Build 2023"

func BuildVersion() string {
	re, err := regexp.Compile(`BuildVersion=(.+?)\n`)
	if err != nil {
		return defaultBuildVer
	}
	match := re.FindStringSubmatch(BuildInfo)

	if len(match) > 0 {
		return match[1]
	}
	return defaultBuildVer
}
