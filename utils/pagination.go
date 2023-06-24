package utils

import (
	"tkspectro/vefeast/app/types/pagination"

	"gorm.io/gorm"
)

// CountMeta counts the total number of records and sets the pagination metadata
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

// TODO: Sorting needs to be completed. This is just a placeholder/simple implementation
// ParseOrder parses the order string and returns a valid order string
func ParseOrder(order string) string {
	if order == "" {
		return "id asc"
	}

	return order
}
