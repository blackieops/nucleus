package nxc

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
