package entity

import (
	"encoding/json"
	"testing"
)

func TestCommonDTOFieldsJSONUsesLowercaseID(t *testing.T) {
	payload, err := json.Marshal(CommonDTOFields{ID: "b1"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded["id"] != "b1" {
		t.Fatalf("expected lowercase id in JSON, got %v", decoded)
	}
	if _, exists := decoded["ID"]; exists {
		t.Fatalf("unexpected uppercase ID key: %v", decoded)
	}
}
