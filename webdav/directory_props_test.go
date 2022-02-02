package webdav

import (
	"time"
	"testing"

	"com.blackieops.nucleus/files"
)

func TestDirectoryPropHandlers(t *testing.T) {
	dir := &files.Directory{
		ID: uint(1234),
		Name: "test",
		FullName: "some/test",
		CreatedAt: time.Date(2022, time.February, 1, 19, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2022, time.February, 1, 19, 0, 0, 0, time.UTC),
	}
	table := []struct {
		key string
		shouldError bool
		expectedTagName string
		expectedValue string
	}{
		{"getlastmodified", false, "d:getlastmodified", "Tue, 01 Feb 2022 19:00:00 UTC"},
		{"getetag", false, "d:getetag", `"Tue, 01 Feb 2022 19:00:00 UTC"`},
		{"permissions", false, "oc:permissions", "RGDNVCK"},
		{"resourcetype", false, "d:resourcetype", "<d:collection/>"},
		{"id", false, "oc:id", "nucleus_1234"},
		{"fileId", false, "oc:fileId", "1234"},
		{"size", false, "oc:size", "0"},
		{"share-types", false, "oc:share-types", ""},

		{"getcontenttype", true, "d:getcontenttype", ""},
		{"getcontentlength", true, "d:getcontentlength", ""},
		{"checksums", true, "oc:checksums", ""},
		{"dDC", true, "oc:dDC", ""},
		{"downloadURL", true, "oc:downloadURL", ""},
	}

	for _, testCase := range table {
		prop, err := DirectoryPropHandlers[testCase.key](dir)
		if err != testCase.shouldError {
			t.Errorf("%s had incorrect error state: %v", testCase.key, testCase.shouldError)
		}
		if prop.XMLName.Local != testCase.expectedTagName {
			t.Errorf("%s had incorrect XML tag name: %s", testCase.key, prop.XMLName.Local)
		}
		if prop.Value != testCase.expectedValue {
			t.Errorf("%s had incorrect value: %s", testCase.key, prop.Value)
		}
	}
}

