package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CRUD[T any] struct {
	DB *sqlx.DB
}

/*
Create returns last inserted ID.
For PostgreSQL, we need "RETURNING id" in the INSERT statement.
*/
func (c *CRUD[T]) Create(ctx context.Context, query string, args ...any) (int64, error) {
	var id int64
	fmt.Println(query, args)
	err := c.DB.QueryRowxContext(ctx, query, args...).Scan(&id)
	return id, err
}

/*
GetOne returns a single struct of type T.
*/
func (c *CRUD[T]) GetOne(ctx context.Context, query string, args ...any) (T, error) {
	var obj T
	err := c.DB.GetContext(ctx, &obj, query, args...)
	return obj, err
}

/*
List returns multiple results of type T.
*/
func (c *CRUD[T]) List(ctx context.Context, query string, args ...any) ([]T, error) {
	var items []T
	err := c.DB.SelectContext(ctx, &items, query, args...)
	return items, err
}

/*
Update returns number of affected rows.
*/
func (c *CRUD[T]) Update(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := c.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

/*
Delete also returns number of affected rows.
*/
func (c *CRUD[T]) Delete(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := c.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
