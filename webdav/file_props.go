package webdav

import (
	"encoding/xml"
	"fmt"
	"time"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/files"
)

// A mapping of property name to a function that takes the file and returns the
// resulting DavProp and a boolean indicating whether or not there is an error.
// Handlers that return `true` for error will be placed empty in the "404 Not
// Found" proplist in the webdav response.
var FilePropHandlers = map[string]func(*files.File) (*DavProp, bool){
	"getlastmodified": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "d:getlastmodified"}, Value: f.UpdatedAt.Format(time.RFC1123)}, false
	},
	"getetag": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "d:getetag"}, Value: `"` + f.Digest + `"`}, false
	},
	"permissions": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:permissions"}, Value: "RGDNVW"}, false
	},
	"resourcetype": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "d:resourcetype"}, Value: ""}, false
	},
	"id": func(f *files.File) (*DavProp, bool) {
		if f.ID == 0 {
			return &DavProp{XMLName: xml.Name{Local: "oc:id"}, Value: ""}, true
		}
		return &DavProp{XMLName: xml.Name{Local: "oc:id"}, Value: "nucleus_" + fmt.Sprint(f.ID)}, false
	},
	"fileId": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:fileId"}, Value: fmt.Sprint(f.ID)}, false
	},
	"size": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:size"}, Value: fmt.Sprint(f.Size)}, false
	},
	"share-types": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:share-types"}, Value: ""}, false
	},
	"downloadURL": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:downloadURL"}, Value: ""}, true
	},
	"dDC": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:dDC"}, Value: ""}, true
	},
	"checksums": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:checksums"}, Value: ""}, true
	},
	"getcontenttype": func(f *files.File) (*DavProp, bool) {
		// TODO: store mimetype on files.File
		return &DavProp{XMLName: xml.Name{Local: "d:getcontenttype"}, Value: "application/octet-stream"}, false
	},
	"getcontentlength": func(f *files.File) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "d:getcontentlength"}, Value: fmt.Sprint(f.Size)}, false
	},
}

func BuildFileProplist(
	file *files.File,
	user *auth.User,
	allowProps []PropfindRequestProp,
) (*DavPropstatList, *DavPropstatList) {
	var okProps []*DavProp
	var errProps []*DavProp
	for _, prop := range allowProps {
		handler := FilePropHandlers[prop.XMLName.Local]
		if handler == nil {
			continue
		}
		result, err := handler(file)
		if err == true {
			errProps = append(errProps, result)
		} else {
			okProps = append(okProps, result)
		}
	}
	return &DavPropstatList{Props: okProps}, &DavPropstatList{Props: errProps}
}
