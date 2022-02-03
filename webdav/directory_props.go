package webdav

import (
	"encoding/xml"
	"fmt"
	"time"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/files"
)

// A mapping of property name to a function that takes the directory and
// returns the resulting DavProp and a boolean indicating whether or not there
// is an error. Handlers that return `true` for error will be placed empty in
// the "404 Not Found" proplist in the webdav response.
var DirectoryPropHandlers = map[string]func(*files.Directory) (*DavProp, bool){
	"getlastmodified": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "d:getlastmodified"}, Value: dir.UpdatedAt.Format(time.RFC1123)}, false
	},
	"getetag": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "d:getetag"}, Value: `"` + dir.UpdatedAt.Format(time.RFC1123) + `"`}, false
	},
	"permissions": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:permissions"}, Value: "RGDNVCK"}, false
	},
	"resourcetype": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "d:resourcetype"}, Value: "<d:collection/>"}, false
	},
	"id": func(dir *files.Directory) (*DavProp, bool) {
		if dir.ID == 0 {
			return &DavProp{XMLName: xml.Name{Local: "oc:id"}, Value: ""}, true
		}
		return &DavProp{XMLName: xml.Name{Local: "oc:id"}, Value: "nucleus_" + fmt.Sprint(dir.ID)}, false
	},
	"fileId": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:fileId"}, Value: fmt.Sprint(dir.ID)}, false
	},
	"size": func(dir *files.Directory) (*DavProp, bool) {
		// TODO: keep track of folder size so we can provide it here
		return &DavProp{XMLName: xml.Name{Local: "oc:size"}, Value: "0"}, false
	},
	"share-types": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:share-types"}, Value: ""}, false
	},
	"downloadURL": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:downloadURL"}, Value: ""}, true
	},
	"dDC": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:dDC"}, Value: ""}, true
	},
	"checksums": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "oc:checksums"}, Value: ""}, true
	},
	"getcontenttype": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "d:getcontenttype"}, Value: ""}, true
	},
	"getcontentlength": func(dir *files.Directory) (*DavProp, bool) {
		return &DavProp{XMLName: xml.Name{Local: "d:getcontentlength"}, Value: ""}, true
	},
}

func BuildDirectoryProplist(
	dir *files.Directory,
	user *auth.User,
	allowProps []PropfindRequestProp,
) (*DavPropstatList, *DavPropstatList) {
	var okProps []*DavProp
	var errProps []*DavProp
	for _, prop := range allowProps {
		handler := DirectoryPropHandlers[prop.XMLName.Local]
		if handler == nil {
			continue
		}
		result, err := handler(dir)
		if err == true {
			errProps = append(errProps, result)
		} else {
			okProps = append(okProps, result)
		}
	}

	return &DavPropstatList{Props: okProps}, &DavPropstatList{Props: errProps}
}
