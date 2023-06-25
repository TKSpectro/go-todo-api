package utils

import (
	"reflect"

	"github.com/TKSpectro/go-todo-api/app/types/pagination"

	"gorm.io/gorm"
)

// CountMeta counts the total number of records and sets the pagination metadata accordingly
func CountMeta(meta *pagination.Meta, query *gorm.DB) {
	count := int64(0)

	query.Count(&count)

	meta.Total = int(count)
	meta.HasNextPage = meta.Offset+meta.Limit < meta.Total
	meta.HasPrevPage = meta.Offset > 0
	meta.TotalPages = int(meta.Total / meta.Limit)
	meta.NextPage = meta.Page + 1
	meta.PrevPage = meta.Page - 1
}

// ParseSearch parses the search string and returns a string that can be used in a gorm query
func ParseSearch(search string) string {
	if search == "" {
		return "%"
	}

	return "%" + search + "%"
}

// SearchWhere returns a string and an array of interfaces that can be used in a gorm query
// The array just needs to be spread into the args of the query like this (query.Where(searchString, searchArray...)
func SearchWhere(search string, model interface{}) (string, []interface{}) {
	amountSearchableFields := 0
	searchString := ""

	t := pointerType(reflect.TypeOf(model))

	firstSearchableField := true
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		xSearchTag := field.Tag.Get("x-search")
		if xSearchTag == "true" {
			amountSearchableFields++

			if firstSearchableField {
				searchString += field.Name + " LIKE ?"

				firstSearchableField = false
			} else {
				searchString += " OR " + field.Name + " LIKE ?"
			}
		}
	}

	searchArray := make([]interface{}, amountSearchableFields)
	for i := 0; i < amountSearchableFields; i++ {
		searchArray[i] = search
	}

	return searchString, searchArray
}
