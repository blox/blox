package regex

const (
	validClusterName = "cluster1"
	validClusterARN  = "arn:aws:ecs:us-east-1:123456789123:cluster/" + validClusterName

	invalidClusterName                 = "cluster1/cluster1"
	invalidClusterARNWithNoName        = "arn:aws:ecs:us-east-1:123456789123:cluster/"
	invalidClusterARNWithInvalidName   = "arn:aws:ecs:us-east-1:123456789123:cluster/" + invalidClusterName
	invalidClusterARNWithInvalidPrefix = "arn/cluster"

	validTaskARN                    = "arn:aws:ecs:us-east-1:123456789012:task/271022c0-f894-4aa2-b063-25bae55088d5"
	invalidTaskARNWithNoID          = "arn:aws:ecs:us-east-1:123456789123:task/"
	invalidTaskARNWithInvalidID     = "arn:aws:ecs:us-east-1:123456789123:task/271022c0-f894-4aa2-b063-25bae55088d5/-"
	invalidTaskARNWithInvalidPrefix = "arn/task"

	validInstanceARN                    = "arn:aws:ecs:us-east-1:123456789123:container-instance/4b6d45ea-a4b4-4269-9d04-3af6ddfdc597"
	invalidInstanceARNWithNoID          = "arn:aws:ecs:us-east-1:123456789123:container-instance/"
	invalidInstanceARNWithInvalidID     = "arn:aws:ecs:us-east-1:123456789123:container-instance/4b6d45ea-a4b4-4269-9d04-3af6ddfdc597/-"
	invalidInstanceARNWithInvalidPrefix = "arn/container-instance"
)
