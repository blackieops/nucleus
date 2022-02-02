package webdav

import (
	"time"

	"com.blackieops.nucleus/auth"
	"com.blackieops.nucleus/files"
)

// A "fake" directory to be used as the user's root directory handle, as the
// "root" directory doesn't really exist but we still need to serialize it in
// some places.
var RootDirectory = files.Directory{
	Name:      "",
	FullName:  "",
	CreatedAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
	UpdatedAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
}

func BuildMultiResponse(
	user *auth.User,
	composite *files.CompositeListing,
	allowProps []PropfindRequestProp,
) *DavMultiResponse {
	var responses []*DavResponse

	if composite.Parent == nil {
		responses = append(responses, BuildDirectoryResponse(&RootDirectory, user, allowProps))
	} else {
		responses = append(responses, BuildDirectoryResponse(composite.Parent, user, allowProps))
	}

	dirResponses := make([]*DavResponse, len(composite.Directories))
	for i, dir := range composite.Directories {
		dirResponses[i] = BuildDirectoryResponse(dir, user, allowProps)
	}
	responses = append(responses, dirResponses...)

	fileResponses := make([]*DavResponse, len(composite.Files))
	for i, f := range composite.Files {
		fileResponses[i] = BuildFileResponse(f, user, allowProps)
	}
	responses = append(responses, fileResponses...)

	return &DavMultiResponse{
		Responses: responses,
		XMLNSNC:   "http://nextcloud.org/ns",
		XMLNSOC:   "http://owncloud.org/ns",
		XMLNSD:    "DAV:",
	}
}

func BuildDirectoryResponse(
	dir *files.Directory,
	user *auth.User,
	allowProps []PropfindRequestProp,
) *DavResponse {
	var dirPath string
	if dir.FullName == "" {
		dirPath = ""
	} else {
		// Ensure we add a trailing slash to all but the root directory path.
		// This is required for Nextcloud to work.
		dirPath = dir.FullName + "/"
	}
	okProps, errProps := BuildDirectoryProplist(dir, user, allowProps)
	return &DavResponse{
		// TODO: pull out path prefix into options somewhere
		Href: "/nextcloud/remote.php/dav/files/" + user.Username + "/" + dirPath,
		Propstats: []*DavPropstat{
			{Props: okProps, Status: "HTTP/1.1 200 OK"},
			{Props: errProps, Status: "HTTP/1.1 404 Not Found"},
		},
	}
}

func BuildFileResponse(
	f *files.File,
	user *auth.User,
	allowProps []PropfindRequestProp,
) *DavResponse {
	okProps, errProps := BuildFileProplist(f, user, allowProps)
	return &DavResponse{
		// TODO: pull out path prefix into options somewhere
		Href: "/nextcloud/remote.php/dav/files/" + user.Username + "/" + f.FullName,
		Propstats: []*DavPropstat{
			{Props: okProps, Status: "HTTP/1.1 200 OK"},
			{Props: errProps, Status: "HTTP/1.1 404 Not Found"},
		},
	}
}
