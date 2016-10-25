package v1

const (
	// 4xx error messages
	instanceNotFoundClientErrMsg = "Instance not found"

	// 5xx error messages
	internalServerErrMsg = "Unexpected internal server error"
	encodingServerErrMsg = "Unexpected server error while encoding response"
	routingServerErrMsg  = "Unexpected server error related to api handler function routing"
)
