package webdav

import (
	"encoding/xml"
	"testing"
	"time"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/files"
)

func TestBuildMultiResponse(t *testing.T) {
	user := &auth.User{
		ID:       123,
		Username: "alice",
	}

	directory := &files.Directory{
		ID:        888,
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Name:      "Notes",
		FullName:  "Notes",
	}

	file := &files.File{
		ID:        999,
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2012, time.November, 12, 10, 30, 0, 0, time.UTC),
		UserID:    123,
		User:      *user,
		Parent:    directory,
		Name:      "test.txt",
		FullName:  "test.txt",
	}

	composite := &files.CompositeListing{
		Parent:      directory,
		Directories: []*files.Directory{},
		Files:       []*files.File{file},
	}

	props := []PropfindRequestProp{
		{XMLName: xml.Name{Local: "getcontentlength"}},
		{XMLName: xml.Name{Local: "getlastmodified"}},
		{XMLName: xml.Name{Local: "permissions"}},
		{XMLName: xml.Name{Local: "checksums"}},
		{XMLName: xml.Name{Local: "fileId"}},
		{XMLName: xml.Name{Local: "id"}},
	}

	r := BuildMultiResponse(user, composite, props)

	// lol sorry
	expected := `<d:multistatus xmlns:nc="http://nextcloud.org/ns" xmlns:oc="http://owncloud.org/ns" xmlns:d="DAV:"><d:response><d:href>/nextcloud/remote.php/dav/files/alice/Notes/</d:href><d:propstat><d:prop><d:getlastmodified>Tue, 10 Nov 2009 23:00:00 UTC</d:getlastmodified><oc:permissions>RGDNVCK</oc:permissions><oc:fileId>888</oc:fileId><oc:id>nucleus_888</oc:id></d:prop><d:status>HTTP/1.1 200 OK</d:status></d:propstat><d:propstat><d:prop><d:getcontentlength></d:getcontentlength><oc:checksums></oc:checksums></d:prop><d:status>HTTP/1.1 404 Not Found</d:status></d:propstat></d:response><d:response><d:href>/nextcloud/remote.php/dav/files/alice/test.txt</d:href><d:propstat><d:prop><d:getcontentlength>0</d:getcontentlength><d:getlastmodified>Mon, 12 Nov 2012 10:30:00 UTC</d:getlastmodified><oc:permissions>RGDNVW</oc:permissions><oc:fileId>999</oc:fileId><oc:id>nucleus_999</oc:id></d:prop><d:status>HTTP/1.1 200 OK</d:status></d:propstat><d:propstat><d:prop><oc:checksums></oc:checksums></d:prop><d:status>HTTP/1.1 404 Not Found</d:status></d:propstat></d:response></d:multistatus>`

	body, err := xml.Marshal(r)
	if err != nil {
		t.Errorf("Error while marshaling xml: %s", err)
	}

	if string(body) != expected {
		t.Errorf("Marshaled XML was unexpected! Got: %s", string(body))
	}
}
