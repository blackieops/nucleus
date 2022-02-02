package webdav

import (
	"encoding/xml"
)

// Represents the whole webdav request as a set of options to modify the
// response. This comes from a variety of sources (headers, body, etc.) and so
// isn't directly unmarshal-able from the request body.
type PropfindOptions struct {
	Depth      int
	Properties []PropfindRequestProp
}

// Represents an XML document that is provded as a request body for a webdav
// PROPFIND that lists the desired response props.
type propfindRequest struct {
	XMLName  xml.Name              `xml:"propfind"`
	Proplist *propfindRequestProps `xml:"prop"`
}

type propfindRequestProps struct {
	XMLName xml.Name              `xml:"prop"`
	Props   []PropfindRequestProp `xml:",any"`
}

type PropfindRequestProp struct {
	XMLName xml.Name
}

func BuildPropfindOptions(depth int, body []byte) (*PropfindOptions, error) {
	var propReq *propfindRequest
	err := xml.Unmarshal(body, &propReq)
	if err != nil {
		return &PropfindOptions{}, err
	}
	propOpts := &PropfindOptions{
		Depth:      depth,
		Properties: propReq.Proplist.Props,
	}
	return propOpts, nil
}
