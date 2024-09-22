package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ccollins476ad/ffhist/dbutil"
	"github.com/ccollins476ad/ffhist/gen"
	"github.com/ccollins476ad/ffhist/model"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	env, err := parseArgs()
	if err != nil {
		fatal(err, true)
	}

	//// Open the places database file.

	var dbPath string
	if env.TempDB {
		// User wants to open a temporary copy rather than the db itself. This
		// is necessary if firefox is currently running (database locked).
		dbPath = gen.TempFilename("ffhist")
		defer os.Remove(dbPath)

		err := gen.CopyFile(env.Path, dbPath)
		if err != nil {
			fatal(err, false)
		}
	} else {
		dbPath = env.Path
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fatal(err, false)
	}
	defer db.Close()

	//// Select requested rows from the moz_places table.

	qs := dbutil.QuerySuffix{
		SortBy:        env.SortBy.String(),
		SortAscending: !env.Reverse,
		Limit:         env.Limit,
	}

	ps, err := model.SelectPlaces(db, qs)
	if err != nil {
		fatal(err, false)
	}

	//// Print the rows to stdout.

	if env.JSON {
		err = printPlacesJSON(env, ps)
	} else {
		err = printPlacesFriendly(env, ps)
	}
	if err != nil {
		fatal(err, false)
	}
}

// fatal optionally prints an error to stderr, optionally prints the rex usage
// text, and terminates with an appropriate status. It prints an error if
// err!=nil. It prints usage text if printUsage==true.
func fatal(err error, printUsage bool) {
	var serr sqlite3.Error
	if errors.As(err, &serr) {
		switch serr.Code {
		case sqlite3.ErrBusy, sqlite3.ErrLocked:
			err = fmt.Errorf("%w: maybe firefox is using the database; specify -t flag", err)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	if printUsage {
		fmt.Fprintln(os.Stderr)
		flag.CommandLine.Usage()
	}

	os.Exit(1)
}

// printPlacesFriendly prints the specified rows using a human-readable table
// format.
func printPlacesFriendly(env *Env, ps []*model.Place) error {
	for _, p := range ps {
		s, err := encodePlaceToFriendlyString(env, p)
		if err != nil {
			return err
		}
		fmt.Println(s)
	}

	return nil
}

// printPlacesJSON prints the specified rows as indended JSON.
func printPlacesJSON(env *Env, ps []*model.Place) error {
	ms := make([]map[string]any, len(ps))

	for i, p := range ps {
		m, err := encodePlaceToMap(env, p)
		if err != nil {
			return err
		}

		ms[i] = m
	}

	s, err := json.MarshalIndent(ms, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(s))
	return nil
}

// encodePlaceToFriendly returns a string representing the given plan using a
// human-readable table row encoding.
func encodePlaceToFriendlyString(env *Env, p *model.Place) (string, error) {
	sb := &strings.Builder{}

	for i, cid := range env.Columns {
		if i != 0 {
			sb.WriteString(" ")
		}

		switch cid {
		case ColumnIDID:
			fmt.Fprintf(sb, "%-10d", p.ID)

		case ColumnIDURL:
			fmt.Fprintf(sb, "%-40s", p.URL)

		case ColumnIDTitle:
			fmt.Fprintf(sb, "%-40s", p.Title)

		case ColumnIDVisitCount:
			fmt.Fprintf(sb, "%-10d", p.VisitCount)

		case ColumnIDFrecency:
			fmt.Fprintf(sb, "%-10d", p.Frecency)

		case ColumnIDLastVisitDate:
			date := p.LastVisitDate.String()
			fmt.Fprintf(sb, "%-40s", date)

		case ColumnIDDescription:
			fmt.Fprintf(sb, "%-40s", p.Description)

		default:
			return "", fmt.Errorf("invalid column id: %v", cid)
		}
	}

	return sb.String(), nil
}

// encodePlaceToMap builds a map from the given place. Keys are column names,
// values are row values.
func encodePlaceToMap(env *Env, p *model.Place) (map[string]any, error) {
	m := make(map[string]any, len(env.Columns))

	for _, cid := range env.Columns {
		val, err := getPlaceField(p, cid)
		if err != nil {
			return nil, err
		}

		m[cid.String()] = val
	}

	return m, nil
}

// getPlaceField returns the given place's data field with the specified column
// ID. That is, it returns one field of a place struct.
func getPlaceField(p *model.Place, cid ColumnID) (any, error) {
	switch cid {
	case ColumnIDID:
		return p.ID, nil

	case ColumnIDURL:
		return p.URL, nil

	case ColumnIDTitle:
		return p.Title, nil

	case ColumnIDVisitCount:
		return p.VisitCount, nil

	case ColumnIDFrecency:
		return p.Frecency, nil

	case ColumnIDLastVisitDate:
		return p.LastVisitDate, nil

	case ColumnIDDescription:
		return p.Description, nil

	default:
		return nil, fmt.Errorf("invalid column id: %v", cid)
	}
}
