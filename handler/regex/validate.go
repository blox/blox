package regex

import "regexp"

func IsClusterARN(clusterARN string) bool {
	validClusterARN := regexp.MustCompile(ClusterARNRegex)
	if validClusterARN.MatchString(clusterARN) {
		return true
	}
	return false
}
