package version

import "fmt"

const Version = "0.1.2"

var (
	Name      string
	GitCommit string

	HumanVersion   = fmt.Sprintf("%s v%s (%s)", Name, Version, GitCommit)
	RunningVersion = fmt.Sprintf("v%s", Version)
)
