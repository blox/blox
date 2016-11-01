package regex

const (
	clusterNameRegexWithoutStart = "[a-zA-Z0-9_]{1,255}$"
	ClusterNameRegex             = "^" + clusterNameRegexWithoutStart
	ClusterARNRegex              = "^(arn:aws:ecs:)([\\-\\w]+):[0-9]{12}:(cluster)/" + clusterNameRegexWithoutStart
	ClusterNameAsARNSuffixRegex  = "/" + clusterNameRegexWithoutStart
	TaskARNRegex                 = "^(arn:aws:ecs):([\\-\\w]+):[0-9]{12}:(task)\\/[\\-\\w]+$"
	InstanceARNRegex             = "^(arn:aws:ecs:)([\\-\\w]+):[0-9]{12}:(container\\-instance)\\/[\\-\\w]+$"
)
