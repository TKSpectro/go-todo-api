package model

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

type Todo struct {
	BaseModel
	Title       zero.String `gorm:"not null" json:"title" x-search:"true" swaggertype:"string" validate:"required,min=1"`
	Description zero.String `gorm:"" json:"description" x-search:"true" swaggertype:"string"`
	Completed   bool        `gorm:"default:false" json:"completed"`
	CompletedAt null.Time   `gorm:"" json:"completedAt" swaggertype:"string" format:"date-time"`

	AccountID uint `gorm:"not null" json:"fkAccountId"`
	// Account   Account
}

func (todo *Todo) New(remote Todo) {
	todo.Title = remote.Title
	todo.Description = remote.Description
	todo.Completed = remote.Completed
	todo.CompletedAt = remote.CompletedAt
}

type todoColumn struct {
	Name           string
	UnmarshalValue func(*Todo, string) error
	MarshalValue   func(*Todo) (string, error)
}

var TodoColumns = []todoColumn{
	{
		Name: "title",
		UnmarshalValue: func(todo *Todo, val string) error {
			todo.Title = zero.StringFrom(val)
			return nil
		},
		MarshalValue: func(todo *Todo) (string, error) {
			return todo.Title.ValueOrZero(), nil
		},
	},
	{
		Name: "description",
		UnmarshalValue: func(todo *Todo, val string) error {
			todo.Description = zero.StringFrom(val)
			return nil
		},
		MarshalValue: func(todo *Todo) (string, error) {
			return todo.Description.ValueOrZero(), nil
		},
	},
	{
		Name: "completed",
		UnmarshalValue: func(todo *Todo, val string) error {
			todo.Completed = val == "true"
			return nil
		},
		MarshalValue: func(todo *Todo) (string, error) {
			return fmt.Sprintf("%t", todo.Completed), nil
		},
	},
	{
		Name: "completed_at",
		UnmarshalValue: func(todo *Todo, val string) error {
			if val == "" {
				todo.CompletedAt.Valid = false
				return nil
			}
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return err
			}
			todo.CompletedAt = null.TimeFrom(t)
			return nil
		},
		MarshalValue: func(todo *Todo) (string, error) {
			if !todo.CompletedAt.Valid {
				return "", nil
			}
			return todo.CompletedAt.ValueOrZero().String(), nil
		},
	},
}

// MarshalRecord encodes the given todo to a record or returns an error.
func (c *Todo) MarshalRecord() ([]string, error) {
	record := make([]string, len(TodoColumns))
	for i, col := range TodoColumns {
		val, err := col.MarshalValue(c)
		if err != nil {
			return nil, err
		}
		record[i] = val
	}
	return record, nil
}

// checkCarColumns returns an error if record doesn't have the right number
// of columns.
func checkColumns(record []string) error {
	if got, want := len(record), len(TodoColumns); got != want {
		return fmt.Errorf("bad number of columns: got=%d want=%d", got, want)
	}
	return nil
}

// UnmarshalRecord decodes the given car record or returns an error.
func (c *Todo) UnmarshalRecord(record []string) error {
	if err := checkColumns(record); err != nil {
		return err
	}
	for i, col := range TodoColumns {
		if err := col.UnmarshalValue(c, record[i]); err != nil {
			return fmt.Errorf("column=%q: %w", col.Name, err)
		}
	}
	return nil
}
