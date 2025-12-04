package client

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRole_UnmarshalJSON_BackwardCompatibility(t *testing.T) {
	t.Run("created_by/updated_by is a string", func(t *testing.T) {
		jsonStr := `{
			"created_by": "user123",
			"updated_by": "user456"
		}`

		var r Role
		if err := json.Unmarshal([]byte(jsonStr), &r); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if r.CreatedBy == nil || r.CreatedBy.ID != "user123" {
			t.Errorf("Expected CreatedBy.ID=user123, got %+v", r.CreatedBy)
		}
		if r.UpdatedBy == nil || r.UpdatedBy.ID != "user456" {
			t.Errorf("Expected UpdatedBy.ID=user456, got %+v", r.UpdatedBy)
		}
	})

	t.Run("created_by/updated_by is an object", func(t *testing.T) {
		jsonObj := `{
			"created_by": {"user_id":"user789","fname":"First A", "lname":"Last A"},
			"updated_by": {"user_id":"user1011","fname":"First B", "lname":"Last B"}
		}`

		expectedCreator := &UserBaseInfo{
			ID:        "user789",
			FirstName: "First A",
			LastName:  "Last A",
		}

		expectedUpdater := &UserBaseInfo{
			ID:        "user1011",
			FirstName: "First B",
			LastName:  "Last B",
		}

		var r2 Role
		if err := json.Unmarshal([]byte(jsonObj), &r2); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !reflect.DeepEqual(r2.CreatedBy, expectedCreator) {
			t.Errorf("Expected CreatedBy=%+v, got %+v", expectedCreator, r2.CreatedBy)
		}

		if !reflect.DeepEqual(r2.UpdatedBy, expectedUpdater) {
			t.Errorf("Expected UpdatedBy=%+v, got %+v", expectedUpdater, r2.UpdatedBy)
		}
	})

	t.Run("created_by/updated_by are null", func(t *testing.T) {
		jsonNull := `{
			"created_by": null,
			"updated_by": null
		}`

		var r3 Role
		if err := json.Unmarshal([]byte(jsonNull), &r3); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if r3.CreatedBy != nil {
			t.Errorf("Expected CreatedBy=nil, got %+v", r3.CreatedBy)
		}
		if r3.UpdatedBy != nil {
			t.Errorf("Expected UpdatedBy=nil, got %+v", r3.UpdatedBy)
		}
	})
}
