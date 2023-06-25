package models

import (
	"reflect"
	"time"

	"github.com/TKSpectro/go-todo-api/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/utils"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updatedAt"`
}

func FindWithMeta(dest interface{}, model interface{}, meta *pagination.Meta, where *gorm.DB) *gorm.DB {
	search, searchArgs := searchWhere(meta.Search, model)

	query := database.DB.Model(model).Where(search, searchArgs...)

	if where != nil {
		query = query.Where(where)
	}

	filters := &meta.Filters
	if filters != nil && len(*filters) > 0 {
		query = query.Where(filtersToQuery(filters))
	}

	countMeta(meta, query)

	return query.
		Offset(meta.Offset).
		Limit(meta.Limit).
		Order(meta.Order).
		Find(dest)
}

func filtersToQuery(filters *[]pagination.FilterEntry) string {
	query := ""

	firstFilter := true

	for _, filter := range *filters {
		if !firstFilter {
			query += " AND "
		} else {
			firstFilter = false
		}

		switch filter.Operator {
		case "n":
			filter.Operator = " IS "
			filter.Value = "NULL"
		case "nn":
			filter.Operator = " IS NOT "
			filter.Value = "NULL"
		case "in":
			filter.Operator = " IN "
			filter.Value = "(" + filter.Value + ")"
		case "nin":
			filter.Operator = " NOT IN "
			filter.Value = "(" + filter.Value + ")"
		case "gt":
			filter.Operator = " > "
			filter.Value = "'" + filter.Value + "'"
		case "gte":
			filter.Operator = " >= "
			filter.Value = "'" + filter.Value + "'"
		case "lt":
			filter.Operator = " < "
			filter.Value = "'" + filter.Value + "'"
		case "lte":
			filter.Operator = " <= "
			filter.Value = "'" + filter.Value + "'"
		case "eq":
			fallthrough // Go does not fall through by default
		default:
			filter.Operator = " = "
			filter.Value = "'" + filter.Value + "'"
		}

		query += filter.Key + filter.Operator + filter.Value
	}

	return query
}

// CountMeta counts the total number of records and sets the pagination metadata accordingly
func countMeta(meta *pagination.Meta, query *gorm.DB) {
	count := int64(0)

	query.Count(&count)

	meta.Total = int(count)
	meta.HasNextPage = meta.Offset+meta.Limit < meta.Total
	meta.HasPrevPage = meta.Offset > 0
	meta.TotalPages = int(meta.Total / meta.Limit)
	meta.NextPage = meta.Page + 1
	meta.PrevPage = meta.Page - 1
}

// SearchWhere returns a string and an array of interfaces that can be used in a gorm query
// The array just needs to be spread into the args of the query like this (query.Where(searchString, searchArray...)
func searchWhere(search string, model interface{}) (string, []interface{}) {
	amountSearchableFields := 0
	searchString := ""

	t := utils.PointerType(reflect.TypeOf(model))

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
