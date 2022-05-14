package nxc

// StatusResponse is the response body structure for the Nextcloud "status"
// endpoint. This is hit periodically by clients to check feature support,
// compatibility, and branding.
type StatusResponse struct {
	Installed            bool   `json:"installed"`
	Maintenance          bool   `json:"maintenance"`
	NeedsDatabaseUpgrade bool   `json:"needsDbUpgrade"`
	Version              string `json:"version"`
	VersionString        string `json:"versionstring"`
	Edition              string `json:"edition"`
	ProductName          string `json:"productname"`
	ExtendedSupport      bool   `json:"extendedSupport"`
}
