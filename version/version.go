package version

import "fmt"

const Version = "0.0.1"

var (
	Name      string
	GitCommit string

	RunningVersion = fmt.Sprintf("v%s", Version)
	HumanVersion   = fmt.Sprintf("%s v%s (%s)", Name, Version, GitCommit)
)
