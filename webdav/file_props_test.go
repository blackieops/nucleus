package webdav

import (
	"time"
	"testing"

	"com.blackieops.nucleus/files"
)

func TestFilePropHandlers(t *testing.T) {
	file := &files.File{
		ID: uint(1234),
		Name: "test",
		FullName: "some/test",
		Size: 69,
		Digest: "71f4376d6551ba7b3363171bf4dc54e41bc18c5e",
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
		{"getetag", false, "d:getetag", `"71f4376d6551ba7b3363171bf4dc54e41bc18c5e"`},
		{"permissions", false, "oc:permissions", "RGDNVW"},
		{"resourcetype", false, "d:resourcetype", ""},
		{"id", false, "oc:id", "nucleus_1234"},
		{"fileId", false, "oc:fileId", "1234"},
		{"size", false, "oc:size", "69"},
		{"getcontenttype", false, "d:getcontenttype", "application/octet-stream"},
		{"getcontentlength", false, "d:getcontentlength", "69"},
		{"share-types", false, "oc:share-types", ""},

		{"checksums", true, "oc:checksums", ""},
		{"dDC", true, "oc:dDC", ""},
		{"downloadURL", true, "oc:downloadURL", ""},
	}

	for _, testCase := range table {
		prop, err := FilePropHandlers[testCase.key](file)
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

