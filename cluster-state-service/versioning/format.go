package versioning

import "fmt"

type Versioner interface {
	Version() (string, error)
}

// PrintVersions prints the version information on stdout as a multi-line
// string. The output will look similar to the following:
//
//    Blox Cluster State Service:
//        Version: 0.0.1
//        Commit: 55347bc
func PrintVersion(extra ...Versioner) {
	cleanliness := ""
	if GitDirty {
		cleanliness = "\tDirty: true\n"
	}

	fmt.Printf(`Blox Cluster State Service:
	Version: %v
	Commit: %v
%v`, Version, GitShortHash, cleanliness)

	for _, versioner := range extra {
		if str, err := versioner.Version(); err == nil {
			fmt.Printf("\t%v\n", str)
		}
	}

}

// String produces a human-readable string showing the agent version.
func String() string {
	ret := "Blox Cluster State Service - v" + Version + " ("
	if GitDirty {
		ret += "*"
	}
	return ret + GitShortHash + ")"
}

func GitHashString() string {
	if GitDirty {
		return "*" + GitShortHash
	}
	return GitShortHash
}
