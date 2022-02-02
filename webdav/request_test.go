package webdav

import (
	"encoding/xml"
	"testing"
)

func TestPropfindRequestUnmarshalsCorrectly(t *testing.T) {
	requestBody := `<d:propfind xmlns:d="DAV:" xmlns:oc="http://owncloud.org/ ns">
		<d:prop>
			<d:resourcetype />
			<d:getlastmodified />
			<d:getcontentlength />
			<d:getetag />
			<oc:size />
			<oc:permissions />
			<oc:checksums />
		</d:prop>
	</d:propfind>`

	var result PropfindRequest
	err := xml.Unmarshal([]byte(requestBody), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal PropfindRequest: %v", err)
		return
	}

	expectedProps := []string{
		"resourcetype",
		"getlastmodified",
		"getcontentlength",
		"getetag",
		"size",
		"permissions",
		"checksums",
	}

	if len(result.Proplist.Props) != 7 {
		t.Errorf("PropfindRequest had the wrong number of props: %d!", len(result.Proplist.Props))
	}

	for i, p := range result.Proplist.Props {
		tagName := p.XMLName.Local
		if tagName != expectedProps[i] {
			t.Errorf("PropfindRequest had unexpected prop: actual=%s expected=%s", tagName, expectedProps[i])
		}
	}
}

func TestBuildPropfindOptions(t *testing.T) {
	fakeBody := `<d:propfind xmlns:d="DAV:" xmlns:oc="http://owncloud.org/ ns">
		<d:prop>
			<d:resourcetype />
			<d:getlastmodified />
			<d:getcontentlength />
			<d:getetag />
			<oc:size />
			<oc:permissions />
			<oc:checksums />
		</d:prop>
	</d:propfind>`
	result, err := BuildPropfindOptions(1, []byte(fakeBody))
	if err != nil {
		t.Errorf("Failed to unmarshal propfind json body: %v", err)
	}
	if len(result.Properties) != 7 {
		t.Errorf("Incorrect number of props for propfind request: %d", len(result.Properties))
	}
}
