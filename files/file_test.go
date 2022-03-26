package files

import (
	"testing"
)

func TestFileSetNames(t *testing.T) {
	// Basic use case
	file := &File{Name: "test.jpg", FullName: "test.jpg"}
	file.SetNames("my cat.jpg")
	if file.FullName != "my cat.jpg" {
		t.Errorf("SetNames did not set the .FullName, value was %v", file.FullName)
	}
	if file.Name != "my cat.jpg" {
		t.Errorf("SetNames did not set the .Name, value was %v", file.Name)
	}

	// Test renaming the Parent but not the filename
	dir := &Directory{Name: "My Documents", FullName: "My Documents"}
	file = &File{Name: "test.jpg", FullName: "some/other/stuff/test.jpg", Parent: dir}
	file.SetNames("test.jpg")
	if file.FullName != "My Documents/test.jpg" {
		t.Errorf("SetNames did not set the .FullName properly, value was %v", file.FullName)
	}
	if file.Name != "test.jpg" {
		t.Errorf("SetNames did something bad to the .Name, value was %v", file.Name)
	}

	// Test renaming the Parent and also the filename
	dir = &Directory{Name: "My Documents", FullName: "My Documents"}
	file = &File{Name: "test.jpg", FullName: "some/other/stuff/test.jpg", Parent: dir}
	file.SetNames("cat.jpg")
	if file.FullName != "My Documents/cat.jpg" {
		t.Errorf("SetNames did not set the .FullName properly, value was %v", file.FullName)
	}
	if file.Name != "cat.jpg" {
		t.Errorf("SetNames did something bad to the .Name, value was %v", file.Name)
	}
}

func TestFileSetDigest(t *testing.T) {
	content := []byte("test contents\n")
	expectedDigest := "40b44f15b4b6690a90792137a03d57c4d2918271"
	file := &File{}
	file.SetDigest(content)
	if file.Digest != expectedDigest {
		t.Errorf("SetDigest did not have expected digest: %v", file.Digest)
	}
}
