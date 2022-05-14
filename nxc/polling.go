package nxc

// PollResponse is the top-level response body for a login v2 poll request that
// has not yet completed.
type PollResponse struct {
	Poll     PollEndpoint `json:"poll"`
	LoginURL string       `json:"login"`
}

// PollEndpoint is exclusively for use inside a PollResponse and contains the
// endpoint to poll for credentials.
type PollEndpoint struct {
	Token       string `json:"token"`
	EndpointURL string `json:"endpoint"`
}

// PollSuccessResponse is the top-level response body for a login v2 poll
// request that has succeeded, and thus contains the auth details for the client
// to use.
type PollSuccessResponse struct {
	Server   string `json:"server"`
	Username string `json:"loginName"`
	Password string `json:"appPassword"`
}
