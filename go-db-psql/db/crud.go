package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"phonax.com/db/metrics"
)

type CRUD[T any] struct {
	DB *sqlx.DB
}

/*
Create returns last inserted ID.
For PostgreSQL, we need "RETURNING id" in the INSERT statement.
*/
func (c *CRUD[T]) Create(ctx context.Context, query string, args ...any) (int64, error) {
	start := time.Now()
	var id int64
	fmt.Println(query, args)
	err := c.DB.QueryRowxContext(ctx, query, args...).Scan(&id)
	if err != nil {
		metrics.RecordQuery("CREATE", time.Since(start), true)
	} else {
		metrics.RecordQuery("CREATE", time.Since(start), false)
	}
	return id, err
}

/*
GetOne returns a single struct of type T.
*/
func (c *CRUD[T]) GetOne(ctx context.Context, query string, args ...any) (T, error) {
	start := time.Now()
	var obj T
	err := c.DB.GetContext(ctx, &obj, query, args...)
	if err != nil {
		metrics.RecordQuery("GET_ONE", time.Since(start), true)
	} else {
		metrics.RecordQuery("GET_ONE", time.Since(start), false)
	}
	return obj, err
}

/*
List returns multiple results of type T.
*/
func (c *CRUD[T]) List(ctx context.Context, query string, args ...any) ([]T, error) {
	start := time.Now()
	var items []T
	err := c.DB.SelectContext(ctx, &items, query, args...)
	if err != nil {
		metrics.RecordQuery("LIST", time.Since(start), true)
	} else {
		metrics.RecordQuery("LIST", time.Since(start), false)
	}
	return items, err
}

/*
Update returns number of affected rows.
*/
func (c *CRUD[T]) Update(ctx context.Context, query string, args ...any) (int64, error) {
	start := time.Now()

	res, err := c.DB.ExecContext(ctx, query, args...)
	fmt.Println(query, args)
	if err != nil {
		return 0, err
	}
	if err != nil {
		metrics.RecordQuery("UPDATE", time.Since(start), true)
	} else {
		metrics.RecordQuery("UPDATE", time.Since(start), false)
	}
	return res.RowsAffected()
}

/*
Delete also returns number of affected rows.
*/
func (c *CRUD[T]) Delete(ctx context.Context, query string, args ...any) (int64, error) {
	start := time.Now()
	res, err := c.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	if err != nil {
		metrics.RecordQuery("DELETE", time.Since(start), true)
	} else {
		metrics.RecordQuery("DELETE", time.Since(start), false)
	}
	return res.RowsAffected()
}
