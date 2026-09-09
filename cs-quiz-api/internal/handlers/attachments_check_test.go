package handlers

import (
	"encoding/json"
	"testing"
)

func TestEncodeDecodeAttachments(t *testing.T) {
	raw, err := encodeAttachments([]string{"  a.pdf ", "", "b.png"})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "a.pdf" || got[1] != "b.png" {
		t.Fatalf("got %#v", got)
	}
	if len(decodeAttachments(nil)) != 0 {
		t.Fatal("nil raw should be empty")
	}
	if len(decodeAttachments([]byte("[]"))) != 0 {
		t.Fatal("empty json array should be empty slice")
	}
}
