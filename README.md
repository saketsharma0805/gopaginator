# paginator

A golang Paginator

## what was the need

For creating a pagination object, that will be passed around in functions, rather than http.Request object or other filters.

## Installation

- Import it in your code

```go
import "github.com/saketsharma0805/paginator"
```

## Quick start

Suppose GetUsers is a function which will return list of users, and we need pagination

```go
func GetUsers (ctx *context.Context, db *sql.DB, p *paginator.Paginator) []*Users {
  stmt := fmt.Sprintf("SELECT id, name, email, is_active, created_at
    from tbl_users
    where 1 and (name like ? or email like ?) and is_active = ?
    order by %s %s
    limit %d offset %d",
    p.OrderBy, p.Order, p.Limit, p.Offset,
  )

  rows, err := db.QueryContext(ctx,stmt, p.Search, p.Search, p.Filter["is_active"])
  ...
 }
```

Pass http request and list of extra filters to Paginator constructor

```go
  p := paginator.NewPagination(r, []string{"is_active"})
  p.ParseRequest()
```

Now we can use p variable to pass in other functions.

```go
  users := GetUsers(ctx, db, p)
```

## Some points to remember

- The purpose of this module was not to validate Query params.
- You can use [Go Playground Validator](https://github.com/go-playground/validator) for validating Query Params.
- Always use p.Search in database searching, and p.Q in templates.
- p.Search encapsulated query with '%' and p.Q is raw text searched by user. p.Search also filters with regex.
- All extra filters are also sanitized with regex, but you have to put other conditions in custom code or may be build a wrapper or something.

## Todo List

- [x] Add Support for filters
- [x] Unit Test
- [ ] Benchmark Test
