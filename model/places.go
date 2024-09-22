package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ccollins476ad/ffhist/dbutil"
)

type Place struct {
	ID            int
	URL           string
	Title         string
	VisitCount    int
	Frecency      int
	LastVisitDate time.Time
	Description   string
}

const (
	PlaceTable   = "moz_places"
	PlaceColumns = "id, url, title, visit_count, frecency, last_visit_date, description"
)

func (p *Place) Scan(rs *sql.Rows) error {
	var (
		id            int
		url           sql.NullString
		title         sql.NullString
		visitCount    int
		frecency      int
		lastVisitDate sql.NullInt64
		description   sql.NullString
	)
	err := rs.Scan(
		&id,
		&url,
		&title,
		&visitCount,
		&frecency,
		&lastVisitDate,
		&description,
	)
	if err != nil {
		return err
	}

	*p = Place{
		ID:            id,
		URL:           url.String,
		Title:         title.String,
		VisitCount:    visitCount,
		Frecency:      frecency,
		LastVisitDate: time.UnixMicro(lastVisitDate.Int64),
		Description:   description.String,
	}
	return nil
}

func SelectPlaces(db *sql.DB, qs dbutil.QuerySuffix) ([]*Place, error) {
	stmt := fmt.Sprintf("select %s from %s", PlaceColumns, PlaceTable)
	suf := qs.String()
	if suf != "" {
		stmt += " " + suf
	}

	rs, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	var ps []*Place
	for rs.Next() {
		p := &Place{}
		err := p.Scan(rs)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	return ps, nil
}
