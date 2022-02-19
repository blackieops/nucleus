package files

import (
	"regexp"
	"testing"

	"com.blackieops.nucleus/auth"
)

var testUser *auth.User = &auth.User{Username: "test"}
var service StorageBackend = &FilesystemBackend{StoragePrefix: "internal/_testdata"}

func TestFilesystemBackendList(t *testing.T) {
	entries := service.List(testUser, "Pictures")
	if len(entries) != 2 {
		t.Errorf("List returned incorrect entry count: %d", len(entries))
	}
	for _, e := range entries {
		matched, err := regexp.MatchString(`(.*)\.png`, e.Name())
		if !matched || err != nil {
			t.Errorf("List returned more than just the expected pngs.")
		}
	}
}

func TestFilesystemBackendStat(t *testing.T) {
	s, err := service.Stat(testUser, "Pictures/Screen Shot 2021-10-11 at 9.10.38 PM.png")
	if err != nil {
		t.Errorf("Stat returned error: %v", err)
	}
	if s.Name() != "Screen Shot 2021-10-11 at 9.10.38 PM.png" {
		t.Errorf("Stat returned wrong filename: %s", s.Name())
	}
	if s.Size() != 15116 {
		t.Errorf("Stat returned wrong filesize: %d", s.Size())
	}
}

func TestFilesystemBackendWriteReadAndDeleteFile(t *testing.T) {
	file := &File{FullName: "test.txt"}
	contents := []byte("this is some content")

	err := service.WriteFile(testUser, file, contents)
	if err != nil {
		t.Errorf("WriteFile could not write: %v", err)
	}

	readBack, err := service.ReadFile(testUser, file)
	if err != nil {
		t.Errorf("Could not read file written by WriteFile: %v", err)
	}
	if string(readBack) != "this is some content" {
		t.Errorf("WriteFile did not write the right content. Got: %s", string(readBack))
	}

	err = service.DeletePath(testUser, file.FullName)
	if err != nil {
		t.Errorf("Failed to delete file: %v", err)
	}
}

func TestFilesystemBackendCreateStatAndDeleteDirectory(t *testing.T) {
	dir := &Directory{FullName: "My Documents/Receipts"}
	err := service.CreateDirectory(testUser, dir)
	if err != nil {
		t.Errorf("CreateDirectory failed: %v", err)
	}
	stat, err := service.Stat(testUser, dir.FullName)
	if err != nil {
		t.Errorf("CreateDirectory didn't make a directory we could stat! %v", err)
	}
	if !stat.IsDir() {
		t.Errorf("CreateDirectory didn't create a directory.")
	}
	err = service.DeletePath(testUser, dir.FullName)
	if err != nil {
		t.Errorf("Couldn't delete directory that CreateDirectory created: %v", err)
	}
}

func TestFilesystemBackendDeletePathWithEmptyPath(t *testing.T) {
	err := service.DeletePath(testUser, "")
	if err == nil {
		t.Errorf("DeletePath allowed destroying entire user folder.")
	}
	err = service.DeletePath(testUser, ".")
	if err == nil {
		t.Errorf("DeletePath allowed destroying entire user folder.")
	}
}
