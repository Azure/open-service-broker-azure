package version

// Values for these are injected by the build
var (
	version string
	commit  string
)

// GetVersion returns the OSBA version. This is either a semantic version
// number or else, in the case of unreleased code, the string "devel".
func GetVersion() string {
	return version
}

// GetCommit returns the git commit SHA for the code that OSBA was built from.
func GetCommit() string {
	return commit
}
