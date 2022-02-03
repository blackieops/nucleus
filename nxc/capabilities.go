package nxc

type OCSCapabilitiesResponse struct {
	Response CapabilitiesResponse `json:"ocs"`
}

type CapabilitiesResponse struct {
	Meta CapabilitiesMeta `json:"meta"`
	Data CapabilitiesData `json:"data"`
}

type CapabilitiesCapabilities struct {
	Theming CapabilitiesTheming `json:"theming"`
	Files   CapabilitiesFiles   `json:"files"`
	Dav     CapabilitiesDav     `json:"dav"`
}

type CapabilitiesData struct {
	Version      CapabilitiesVersion      `json:"version"`
	Capabilities CapabilitiesCapabilities `json:"capabilities"`
}

type CapabilitiesMeta struct{}

type CapabilitiesFiles struct {
	BigFileChunking bool `json:"bigfilechunking"`
}

type CapabilitiesDav struct {
	ChunkingAPIVersion string `json:"chunking"`
}

type CapabilitiesVersion struct {
	Major           int    `json:"major"`
	Minor           int    `json:"minor"`
	Micro           int    `json:"micro"`
	String          string `json:"string"`
	Edition         string `json:"edition"`
	ExtendedSupport bool   `json:"extendedSupport"`
}

type CapabilitiesTheming struct {
	Name string `json:"name"`
}

func BuildCapabilitiesResponse() *OCSCapabilitiesResponse {
	return &OCSCapabilitiesResponse{
		Response: CapabilitiesResponse{
			Meta: CapabilitiesMeta{},
			Data: CapabilitiesData{
				Version: CapabilitiesVersion{
					Major:           22,
					Minor:           2,
					Micro:           3,
					String:          "22.2.3",
					ExtendedSupport: false,
				},
				Capabilities: CapabilitiesCapabilities{
					Files:   CapabilitiesFiles{BigFileChunking: true},
					Theming: CapabilitiesTheming{Name: "Nucleus"},
					Dav:     CapabilitiesDav{ChunkingAPIVersion: "1.0"},
				},
			},
		},
	}
}
