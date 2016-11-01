package regex

import "regexp"

func IsClusterName(clusterName string) bool {
	validClusterName := regexp.MustCompile(ClusterNameRegex)
	if validClusterName.MatchString(clusterName) {
		return true
	}
	return false
}

func IsClusterARN(clusterARN string) bool {
	validClusterARN := regexp.MustCompile(ClusterARNRegex)
	if validClusterARN.MatchString(clusterARN) {
		return true
	}
	return false
}

func IsTaskARN(taskARN string) bool {
	validTaskARN := regexp.MustCompile(TaskARNRegex)
	if validTaskARN.MatchString(taskARN) {
		return true
	}
	return false
}

func IsInstanceARN(instanceARN string) bool {
	validInstanceARN := regexp.MustCompile(InstanceARNRegex)
	if validInstanceARN.MatchString(instanceARN) {
		return true
	}
	return false
}
