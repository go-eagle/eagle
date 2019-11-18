package model

import (
	"fmt"
	"math"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type BaseModel struct {
	Id        uint64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"-"`
	DeletedAt *time.Time `gorm:"column:deleted_at" sql:"index" json:"-"`
}

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

type Validator interface {
	Validate() error
}

type ValidateChecker interface {
	IsValidated() bool
}

type ValidatedObject interface {
	ValidateChecker
	Validator
}

var ValidatedDefaultValidateChecker = DefaultValidateChecker{isValidated: true}

type DefaultValidateChecker struct {
	isValidated bool `json:"-"`
}

func (c DefaultValidateChecker) IsValidated() bool {
	return c.isValidated
}

type ListResponse struct {
	TotalCount uint64      `json:"total_count"`
	HasMore    int         `json:"has_more"`
	PageKey    string      `json:"page_key"`
	PageValue  int         `json:"page_value"`
	Items      interface{} `json:"items"`
}

type PaginatedList struct {
	CurrentPage int
	NumItem     int
	TotalCount  int
	Items       interface{}
}

func (list PaginatedList) MarshalJSON() ([]byte, error) {
	m := struct {
		TotalCount  int         `json:"total_count"`
		TotalPage   int         `json:"total_page"`
		CurrentPage int         `json:"current_page"`
		HasNextPage bool        `json:"has_next_page"`
		Items       interface{} `json:"items"`
	}{
		TotalCount:  list.TotalCount,
		CurrentPage: list.CurrentPage,
		Items:       list.Items,
	}

	totalPageBeforeCeil := list.TotalCount / list.NumItem
	m.TotalPage = int(math.Ceil(float64(totalPageBeforeCeil)))
	if m.TotalPage == 0 {
		m.TotalPage = 1
	}
	m.HasNextPage = list.CurrentPage < m.TotalPage

	return jsoniter.Marshal(m)
}

// sql build where
// see: https://github.com/jinzhu/gorm/issues/2055
func WhereBuild(where map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}

		if whereSQL != "" {
			whereSQL += " AND "
		}

		fmt.Println(strings.Join(ks, ","))
		switch len(ks) {
		case 1:
			//fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				fmt.Println()
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				vals = append(vals, v)
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)
			}
			break
		}
	}
	return
}
