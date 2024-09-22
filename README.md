# ffhist

ffhist is a tool that prints the contents of a firefox "places" database.

## Usage

```
Usage: ffhist [option]... <db_path>
Prints the contents of a firefox "places" database.
  -c string
        comma-separate list of column names to display (default "last_visit_date,url")
  -j    encode output as json
  -n int
        max rows to include in output (0 for no limit)
  -r    reverse sort order
  -s string
        name of field to sort by (default "last_visit_date")
  -t    copy database to temp file before opening; required if firefox is running

Columns:
  id
  url
  title
  visit_count
  frecency
  last_visit_date
  description
```

## Examples

Print last visit date, url, and title. Sort by last visit date.

```
ffhist -t -c last_visit_date,url,title -s last_visit_date /home/me/.mozilla/firefox/fns5l92s.default-release/places.sqlite
```
