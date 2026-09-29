package tests

import (
	"campus-lost-found-api/app/model"
	"campus-lost-found-api/app/service"
	"testing"
)

func TestApplyPatchOnlyChangesSentFields(t *testing.T) {
	item := model.Item{Title: "A", Description: "desc", Category: "book", Location: "lab", Status: "lost"}
	newTitle := "B"
	service.ApplyPatch(&item, model.PatchItemRequest{Title: &newTitle})
	if item.Title != "B" || item.Description != "desc" {
		t.Fatal("field yang tidak dikirim ikut berubah")
	}
}
func TestCanModifyOwner(t *testing.T) {
	if !service.CanModifyItem("user", 5, 6) {
	} else {
		t.Fatal("user lain tidak boleh mengubah item")
	}
	if !service.CanModifyItem("user", 5, 5) {
		t.Fatal("pemilik harus boleh mengubah item")
	}
}
func TestCanModifyStaff(t *testing.T) {
	if !service.CanModifyItem("staff", 5, 6) {
		t.Fatal("staff harus boleh mengubah item")
	}
}
