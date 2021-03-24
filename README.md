# gopaginator
A golang Paginator

## what was the need
For creating a pagination object, that will be passed around in functions, rather than http.Request object or other filters. 

## Installation
- Import it in your code
``` import "github.com/saketsharma0805/gopaginator"```

## Quick start
```
  # Suppose GetUsers is a function which will return list of users, and we need pagination
  func GetUsers (ctx *context.Context, db *sql.DB, p *gopaginator.Paginator) []*Users {
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

``` 
  # first parameter: here r is of type *http.Request
  # second parameter: []string{} is for the list of extra filters, other than: limit, page, q.
  # in this example, it is "is_active"
  p := gopaginator.NewPaginator(r, []string{"is_active"})
  p.ParseRequest()
```
```
  # Now we can use p variable to pass in other functions.
  users := GetUsers(ctx, db, p)
```

## Todo List
- [X] Add Support for filters
- [X] Unit Test
- [ ] Benchmark Test

