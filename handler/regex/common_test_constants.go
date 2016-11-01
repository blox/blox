package regex

const (
	clusterName     = "cluster1"
	validClusterARN = "arn:aws:ecs:us-east-1:123456789123:cluster/" + clusterName

	invalidClusterARNWithNoName        = "arn:aws:ecs:us-east-1:123456789123:cluster/"
	invalidClusterARNWithInvalidName   = "arn:aws:ecs:us-east-1:123456789123:cluster/cluster1/cluster1"
	invalidClusterARNWithInvalidPrefix = "arn/cluster"
)
