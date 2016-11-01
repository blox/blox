package regex

const (
	ClusterNameRegex            = "[a-zA-Z0-9_]{1,255}"
	ClusterARNRegex             = "(arn:aws:ecs:)([\\-\\w]+):[0-9]{12}:(cluster)/" + ClusterNameRegex
	ClusterNameAsARNSuffixRegex = "/" + ClusterNameRegex
	TaskARNRegex                = "(arn:aws:ecs):([\\-\\w]+):[0-9]{12}:(task)\\/[\\-\\w]+"
	InstanceARNRegex            = "(arn:aws:ecs:)([\\-\\w]+):[0-9]{12}:(container\\-instance)\\/[\\-\\w]+"
)
