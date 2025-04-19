# remitly_hw

### Parsing

[ref1](https://shaileshb.hashnode.dev/go-csv-parsing) \
[ref2](https://gosamples.dev/read-csv/) - example of read-csv \
[ref3](https://pkg.go.dev/encoding/csv#section-sourcefiles) - refer to csv package\
[ref4](https://github.com/gocarina/gocsv/blob/78e41c74b4b1/examples/full/main.go)

[Read csv file into a slice of record structure](https://shaileshb.hashnode.dev/go-csv-parsing)

### Check duplicate

[ref_1](https://stackoverflow.com/questions/39086976/golang-csv-remove-duplicate-if-matching-column-values)

### Set up database

[ref_1](https://go.dev/doc/tutorial/database-access) - Introduction to connect MySQL database in GO

- change port number in `/usr/local/mysql/bin 
╰─$ cat  mysql_config `

[ref_2](https://go.dev/doc/database/) - Accessing relational database

The `database/sql` package you’ll be using includes types and functions for connecting to databases, executing transactions, canceling an operation in progress, and more

1. Install driver `go get -u github.com/go-sql-driver/mysql`
2. Go MySQL Driver is an implementation of Go's `database/sql/driver`interface. You only need to import the driver and can use the full `database/sql` API then.
