package nxc

type PollResponse struct {
	Poll     PollEndpoint `json:"poll"`
	LoginURL string       `json:"login"`
}

type PollEndpoint struct {
	Token       string `json:"token"`
	EndpointURL string `json:"endpoint"`
}

type PollSuccessResponse struct {
	Server   string `json:"server"`
	Username string `json:"loginName"`
	Password string `json:"appPassword"`
}
