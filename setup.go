package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ColumnID int

const (
	ColumnIDID ColumnID = iota
	ColumnIDURL
	ColumnIDTitle
	ColumnIDVisitCount
	ColumnIDFrecency
	ColumnIDLastVisitDate
	ColumnIDDescription
)

var ColumnNames = []string{
	ColumnIDID:            "id",
	ColumnIDURL:           "url",
	ColumnIDTitle:         "title",
	ColumnIDVisitCount:    "visit_count",
	ColumnIDFrecency:      "frecency",
	ColumnIDLastVisitDate: "last_visit_date",
	ColumnIDDescription:   "description",
}

var ColumnNameIDMap = map[string]ColumnID{}

func init() {
	for id, name := range ColumnNames {
		ColumnNameIDMap[name] = ColumnID(id)
	}
}

func (cid ColumnID) String() string {
	return ColumnNames[cid]
}

type Env struct {
	Path    string
	TempDB  bool
	Columns []ColumnID
	SortBy  ColumnID
	Reverse bool
	JSON    bool
	Limit   int
}

func parseArgs() (*Env, error) {
	tempDB := flag.Bool("t", false, "copy database to temp file before opening; required if firefox is running")
	columns := flag.String("c", "last_visit_date,url", "comma-separate list of column names to display")
	sortBy := flag.String("s", "last_visit_date", "name of field to sort by")
	reverse := flag.Bool("r", false, "reverse sort order")
	jsonEncoding := flag.Bool("j", false, "encode output as json")
	limit := flag.Int("n", 0, "max rows to include in output (0 for no limit)")

	flag.Usage = usage
	flag.Parse()

	colIDs, err := parseColumnList(*columns)
	if err != nil {
		return nil, fmt.Errorf("failed to parse column list: %w", err)
	}

	sortID, err := columnNameToID(*sortBy)
	if err != nil {
		return nil, fmt.Errorf("invalid sort specifier: %w", err)
	}

	// Required argument specifies database path.
	if len(flag.Args()) == 0 {
		return nil, fmt.Errorf("missing required argument: db_path")
	}
	dbPath := flag.Args()[0]

	return &Env{
		Path:    dbPath,
		TempDB:  *tempDB,
		Columns: colIDs,
		SortBy:  sortID,
		Reverse: *reverse,
		JSON:    *jsonEncoding,
		Limit:   *limit,
	}, nil
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [option]... <db_path>\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(flag.CommandLine.Output(), `Prints the contents of a firefox "places" database.\n\n`)
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Columns:\n")
	for _, c := range ColumnNames {
		fmt.Fprintf(flag.CommandLine.Output(), "  %s\n", c)
	}
}

func parseColumnList(names string) ([]ColumnID, error) {
	elems := strings.Split(names, ",")

	var cids []ColumnID
	for _, elem := range elems {
		cid, err := columnNameToID(elem)
		if err != nil {
			return nil, err
		}
		cids = append(cids, cid)
	}

	return cids, nil
}

func columnNameToID(name string) (ColumnID, error) {
	id, ok := ColumnNameIDMap[name]
	if !ok {
		return 0, fmt.Errorf("invalid column name: %s", name)
	}

	return id, nil
}
