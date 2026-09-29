package service

import "campus-lost-found-api/app/model"

func ApplyPatch(item *model.Item, p model.PatchItemRequest) {
	if p.Title != nil {
		item.Title = *p.Title
	}
	if p.Description != nil {
		item.Description = *p.Description
	}
	if p.Category != nil {
		item.Category = *p.Category
	}
	if p.Location != nil {
		item.Location = *p.Location
	}
	if p.Status != nil {
		item.Status = *p.Status
	}
}
func CanModifyItem(role string, ownerID, userID int64) bool {
	return role == "admin" || role == "staff" || ownerID == userID
}
