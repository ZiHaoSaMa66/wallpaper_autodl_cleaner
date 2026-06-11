package model

import (
	"encoding/json"
	"testing"
)

func TestGetUserFilesResponse_Unmarshal_publishedfiledetails(t *testing.T) {
	raw := `{"response":{"total":3,"publishedfiledetails":[{"publishedfileid":"123"},{"publishedfileid":"456"},{"publishedfileid":"789"}]}}`
	var resp GetUserFilesResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Response.Total != 3 {
		t.Fatalf("expected total 3, got %d", resp.Response.Total)
	}
	if len(resp.Response.Files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(resp.Response.Files))
	}
	if resp.Response.Files[0].PublishedFileID != "123" {
		t.Fatalf("expected first id 123, got %s", resp.Response.Files[0].PublishedFileID)
	}
}

func TestGetUserFilesResponse_Unmarshal_wrongFieldName_fails(t *testing.T) {
	raw := `{"response":{"total":3,"files":[{"publishedfileid":"123"}]}}`
	var resp GetUserFilesResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Response.Files) != 0 {
		t.Fatal("expected 0 files when using wrong field name 'files'")
	}
}
