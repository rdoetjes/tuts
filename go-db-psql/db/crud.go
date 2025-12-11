package db

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"phonax.com/db/metrics"
)

type CRUD[T any] struct {
	DB *sqlx.DB
}

func CheckConnection(db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}

/*Handle errors and record metrics for all CRUD operations.*/
func (c *CRUD[T]) handleErrorAndMetrics(start time.Time, ACTION string, err error) {
	if err != nil {
		log.Println("ERROR ", ACTION, ":", err)
		metrics.RecordQuery(ACTION, time.Since(start), true)
	} else {
		metrics.RecordQuery(ACTION, time.Since(start), false)
	}
}

/*
Create returns last inserted ID.
For PostgreSQL, we need "RETURNING id" in the INSERT statement.
*/
func (c *CRUD[T]) Create(ctx context.Context, query string, args ...any) (uint64, error) {
	start := time.Now()
	const ACTION = "CREATE"
	var id uint64

	CheckConnection(c.DB)

	err := c.DB.QueryRowxContext(ctx, query, args...).Scan(&id)
	c.handleErrorAndMetrics(start, ACTION, err)
	return id, err
}

/*
GetOne returns a single struct of type T.
*/
func (c *CRUD[T]) GetOne(ctx context.Context, query string, args ...any) (T, error) {
	start := time.Now()
	const ACTION = "GET_ONE"
	var obj T

	CheckConnection(c.DB)

	err := c.DB.GetContext(ctx, &obj, query, args...)
	c.handleErrorAndMetrics(start, ACTION, err)
	return obj, err
}

/*
List returns multiple results of type T.
*/
func (c *CRUD[T]) List(ctx context.Context, query string, args ...any) ([]T, error) {
	start := time.Now()
	const ACTION = "LIST"
	var items []T

	CheckConnection(c.DB)

	err := c.DB.SelectContext(ctx, &items, query, args...)
	c.handleErrorAndMetrics(start, ACTION, err)
	return items, err
}

/*
Update returns number of affected rows.
*/
func (c *CRUD[T]) Update(ctx context.Context, query string, args ...any) (int64, error) {
	start := time.Now()
	const ACTION = "UPDATE"

	CheckConnection(c.DB)

	res, err := c.DB.ExecContext(ctx, query, args...)
	if err != nil {
		metrics.RecordQuery(ACTION, time.Since(start), true)
		log.Println("ERROR ", ACTION, ":", err)
		return 0, err
	}
	metrics.RecordQuery(ACTION, time.Since(start), false)

	return res.RowsAffected()
}

/*
Delete also returns number of affected rows.
*/
func (c *CRUD[T]) Delete(ctx context.Context, query string, args ...any) (int64, error) {
	start := time.Now()
	const ACTION = "DELETE"

	CheckConnection(c.DB)

	res, err := c.DB.ExecContext(ctx, query, args...)
	if err != nil {
		metrics.RecordQuery(ACTION, time.Since(start), true)
		log.Println("ERROR ", ACTION, ":", err)
		return 0, err
	}
	metrics.RecordQuery(ACTION, time.Since(start), false)

	return res.RowsAffected()
}
