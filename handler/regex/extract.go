package regex

import (
	"errors"
	"regexp"
)

func GetClusterNameFromARN(clusterARN string) (string, error) {
	if len(clusterARN) == 0 {
		return "", errors.New("Cluster ARN cannot be empty")
	}

	if !IsClusterARN(clusterARN) {
		return "", errors.New("Invalid cluster ARN")
	}

	re := regexp.MustCompile(ClusterNameAsARNSuffixRegex)
	matchedStrs := re.FindStringSubmatch(clusterARN)
	if len(matchedStrs) != 1 {
		return "", errors.New("Unable to extract cluster name from cluster ARN")
	}

	// matchedStrs[0]=/clusterName. Strip off "/" in the beginning.
	clusterName := matchedStrs[0][1:]
	return clusterName, nil
}
