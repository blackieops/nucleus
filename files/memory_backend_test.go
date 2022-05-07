package files

import (
	"go.b8s.dev/nucleus/auth"
	"testing"
)

func TestMemoryStorageBackendWriteFile(t *testing.T) {
	u := &auth.User{}
	b := &MemoryStorageBackend{}
	f := &File{Name: "thing.txt", FullName: "thing.txt", Size: 4}

	b.WriteFile(u, f, []byte("butt"))

	if string(b.entries["thing.txt"]) != "butt" {
		t.Errorf("WriteFile did not have correct file contents: %v", b.entries["thing.txt"])
	}
}

func TestMemoryStorageBackendReadFile(t *testing.T) {
	u := &auth.User{}
	b := &MemoryStorageBackend{
		entries: map[string][]byte{
			"thing.txt": []byte("works"),
		},
	}
	f := &File{Name: "thing.txt"}

	b.ReadFile(u, f)

	if string(b.entries["thing.txt"]) != "works" {
		t.Errorf("ReadFile did not have correct file contents: %v", b.entries["thing.txt"])
	}
}

func TestMemoryStorageBackendList(t *testing.T) {
	u := &auth.User{}
	b := &MemoryStorageBackend{
		entries: map[string][]byte{
			"thing.txt":      []byte("wrong"),
			"yoink/butt.rtf": []byte("idk"),
			"yoink/test.txt": []byte("correct"),
		},
	}

	items, _ := b.List(u, "yoink")

	if len(items) != 2 {
		t.Errorf("List did not list correct files: %v", items)
	}
}

func TestMemoryStorageBackendStat(t *testing.T) {
	u := &auth.User{}
	b := &MemoryStorageBackend{
		entries: map[string][]byte{
			"thing.txt": []byte("yep"),
		},
	}

	stat, _ := b.Stat(u, "thing.txt")

	if stat.Size() != 3 {
		t.Errorf("Stat did not have correct file size: %v", stat.Size())
	}
	if stat.Name() != "thing.txt" {
		t.Errorf("Stat did not have correct file name: %v", stat.Name())
	}
	if stat.IsDir() != false {
		t.Errorf("Stat did not have correct IsDir value: %v", stat.IsDir())
	}
}

func TestMemoryStorageBackendDelete(t *testing.T) {
	u := &auth.User{}
	b := &MemoryStorageBackend{
		entries: map[string][]byte{
			"thing.txt":      []byte("wrong"),
			"yoink/butt.rtf": []byte("idk"),
			"yoink/test.txt": []byte("correct"),
		},
	}

	b.DeletePath(u, "yoink")

	if len(b.entries) != 1 {
		t.Errorf("DeletePath did not delete enough files: still have %v", len(b.entries))
	}
}

func TestMemoryStorageBackendRename(t *testing.T) {
	u := &auth.User{}
	b := &MemoryStorageBackend{
		entries: map[string][]byte{
			"thing.txt":      []byte("wrong"),
			"yoink/butt.rtf": []byte("idk"),
			"yoink/test.txt": []byte("correct"),
		},
	}

	b.RenamePath(u, "yoink", "boink")

	if len(b.entries) != 3 {
		t.Errorf("RenamePath did something wrong, now we have %v files", len(b.entries))
	}
	if string(b.entries["yoink/test.txt"]) != "correct" {
		t.Errorf("RenamePath did not rename test.txt correctly! content: %v", string(b.entries["yoin/test.txt"]))
	}
}
