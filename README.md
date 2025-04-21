# remitly_hw

## Environment for programming

Database: MYSQL 8.0.35
Language : go version go1.22.5 darwin/amd64
Endpoint API check : Postman

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
[ref_3](https://golangbot.com/mysql-create-table-insert-row/) - Modularize DB connection and create table

The `database/sql` package you’ll be using includes types and functions for connecting to databases, executing transactions, canceling an operation in progress, and more

1. Install driver `go get -u github.com/go-sql-driver/mysql`
2. Go MySQL Driver is an implementation of Go's `database/sql/driver`interface. You only need to import the driver and can use the full `database/sql` API then.
3. Set `DBUSER` and `DBPASS` to login your database.

```shell
$ export DBUSER=username
$ export DBPASS=password
```

### Misellinous

- a "chan" (short for channel) is a communication mechanism that allows goroutines (lightweight threads) to communicate with each other and synchronize their execution.
-

### RESTful API

[ref-1](https://go.dev/doc/tutorial/web-service-gin) - Introduction RESTful API with GO and Gin

- Need to install `go get -u github.com/gin-gonic/gin` to use gin web framework
- [ref-2](https://gin-gonic.com/en/docs/quickstart/) - gin
  [ref-3](https://go.dev/doc/tutorial/web-service-gin#write-the-code) - Step by step to RESTful API for GO
- `gin.Context` is the most important part of Gin. It carries request details, validates and serializes JSON, and more. (Despite the similar name, this is different from Go’s built-in `context` package.)

- Call `Context.IndentedJSON` to serialize the struct into JSON and add it to the response.

- [parameters in path](https://gin-gonic.com/en/docs/examples/param-in-path/)
