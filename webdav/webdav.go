package webdav

import (
	"net/http"
)

type Config struct {
	Storage FileStorage
	Index   FileIndex
}

// FileStorage provides an interface to connect with some mass-storage backend,
// such as a filesystem or remote object storage API, etc. Doesn't matter where
// the data is stored, as long as it conforms to the basic interface here.
type FileStorage interface {
}

// FileIndex represents a database of some kind that storese metadata about all
// the files stored in the FileStorage-compliant backend, as well as the
// Properties for those files.
type FileIndex interface {
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement webdav lol
}
