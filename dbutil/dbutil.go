package dbutil

import (
	"fmt"
	"strings"
)

//type SQLScanner interface {
//Scan(rs *sql.Rows) error
//}

type QuerySuffix struct {
	SortBy        string // "" for no sort
	SortAscending bool
	Limit         int // <=0 for no limit
}

func (q *QuerySuffix) String() string {
	var parts []string

	if q.SortBy != "" {
		ss := SortSQL(q.SortBy, q.SortAscending)
		parts = append(parts, ss)
	}

	if q.Limit > 0 {
		ls, err := LimitSQL(q.Limit)
		if err != nil {
			panic(err)
		}
		parts = append(parts, ls)
	}

	return strings.Join(parts, " ")
}

func SortSQL(column string, ascending bool) string {
	var order string
	if ascending {
		order = "asc"
	} else {
		order = "desc"
	}
	return fmt.Sprintf("order by %s %s", column, order)
}

func LimitSQL(limit int) (string, error) {
	if limit <= 0 {
		return "", fmt.Errorf("invalid sql limit: have=%d want>0", limit)
	}

	return fmt.Sprintf("limit %d", limit), nil
}
