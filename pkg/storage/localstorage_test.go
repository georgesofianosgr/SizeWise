package storage

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestLocal_WriteEntry(t *testing.T) {
	// TODO: Write to tmp file
	t.Fatal("not implemented")
}

func TestLocal_ReadEntry(t *testing.T) {
	basePath := "/tmp/"
	const entryContent = "test passed"
	storage := NewLocal("host", basePath)
	//nolint:gosec
	randomFileName := strconv.FormatInt(rand.Int63n(1000), 10) + ".txt"

	err := storage.WriteEntry(randomFileName, strings.NewReader(entryContent))
	if err != nil {
		t.Fatalf("got error %q", err)
	}
	content := bytes.Buffer{}
	err = storage.ReadEntry(randomFileName, &content)
	if err != nil {
		t.Fatalf("got error %q", err)
	}
	if content.String() != entryContent {
		t.Fatalf("expected content to be 'test', got %q", content)
	}
}
