package handler

import (
	"reflect"

	"github.com/TKSpectro/go-todo-api/pkg/app/service"
	"github.com/TKSpectro/go-todo-api/pkg/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/utils"
	"gorm.io/gorm"
)

type Handler struct {
	accountService service.IAccountService
	todoService    service.ITodoService
	db             *gorm.DB
	validator      *Validator
}

func NewHandler(db *gorm.DB, as service.IAccountService, ts service.ITodoService) *Handler {
	v := NewValidator()

	return &Handler{
		accountService: as,
		todoService:    ts,
		db:             db,
		validator:      v,
	}
}

func (h *Handler) FindWithMeta(dest interface{}, model interface{}, meta *pagination.Meta, where *gorm.DB) *gorm.DB {
	search, searchArgs := searchWhere(meta.Search, model)

	query := h.db.Model(model).Where(search, searchArgs...)

	if where != nil {
		query = query.Where(where)
	}

	filters := &meta.Filters
	if filters != nil && len(*filters) > 0 {
		query = query.Where(filtersToQuery(filters))
	}

	query = query.Offset(meta.Offset).Limit(meta.Limit)

	orders := &meta.Order
	if orders != nil && len(*orders) > 0 {
		query = query.Order(ordersToQuery(orders))
	}

	countMeta(meta, query)

	return query.Find(dest)
}

func ordersToQuery(orders *[]pagination.OrderEntry) string {
	query := ""

	firstOrder := true

	for _, order := range *orders {
		if !firstOrder {
			query += ", "
		} else {
			firstOrder = false
		}

		query += order.Key + " " + order.Direction
	}

	return query
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
