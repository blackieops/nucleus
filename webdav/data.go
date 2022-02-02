package webdav

import (
	"encoding/xml"
)

type DavMultiResponse struct {
	XMLName   xml.Name `xml:"d:multistatus"`
	Responses []*DavResponse
	XMLNSNC   string `xml:"xmlns:nc,attr"`
	XMLNSOC   string `xml:"xmlns:oc,attr"`
	XMLNSD    string `xml:"xmlns:d,attr"`
}

type DavResponse struct {
	XMLName   xml.Name       `xml:"d:response"`
	Href      string         `xml:"d:href"`
	Propstats []*DavPropstat `xml:"d:propstat"`
}

type DavPropstat struct {
	Props  *DavPropstatList `xml:"d:prop"`
	Status DavStatus        `xml:"d:status"`
}

type DavPropstatList struct {
	Props []*DavProp
}

type DavProp struct {
	XMLName xml.Name
	Value   string `xml:",innerxml"`
}

type DavStatus string
