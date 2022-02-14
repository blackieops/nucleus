package files

import (
	"testing"
)

func TestDirectorySetNames(t *testing.T) {
	// Basic use case
	file := &Directory{Name: "Pix", FullName: "Pix"}
	file.SetNames("Photos")
	if file.FullName != "Photos" {
		t.Errorf("SetNames did not set the .FullName, value was %v", file.FullName)
	}
	if file.Name != "Photos" {
		t.Errorf("SetNames did not set the .Name, value was %v", file.Name)
	}

	// Test renaming the Parent but not the Directory
	dir := &Directory{Name: "My Documents", FullName: "My Documents"}
	file = &Directory{Name: "tests", FullName: "some/other/stuff/tests", Parent: dir}
	file.SetNames("tests")
	if file.FullName != "My Documents/tests" {
		t.Errorf("SetNames did not set the .FullName properly, value was %v", file.FullName)
	}
	if file.Name != "tests" {
		t.Errorf("SetNames did something bad to the .Name, value was %v", file.Name)
	}

	// Test renaming the Parent and also the Directory's name
	dir = &Directory{Name: "My Documents", FullName: "My Documents"}
	file = &Directory{Name: "tests", FullName: "some/other/stuff/tests", Parent: dir}
	file.SetNames("Photos")
	if file.FullName != "My Documents/Photos" {
		t.Errorf("SetNames did not set the .FullName properly, value was %v", file.FullName)
	}
	if file.Name != "Photos" {
		t.Errorf("SetNames did something bad to the .Name, value was %v", file.Name)
	}
}
