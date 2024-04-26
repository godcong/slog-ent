package entslog

import (
	"context"
	"fmt"
	"log/slog"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
)

// SlogDriver is a init that logs all init operations.
type SlogDriver struct {
	dialect.Driver // underlying init.
	Handler        // log function. defaults to slog.Default()
}

// New gets a init and an optional logging function, and returns
// a new slog-init that prints all outgoing operations.
func New(logger *slog.Logger, opts ...func(*Option) *Option) func(dialect.Driver) dialect.Driver {
	op := settings(opts...)

	handle := op.handle(logger)
	return func(dri dialect.Driver) dialect.Driver {
		return &SlogDriver{Driver: dri, Handler: handle.initWith(slog.String("database", "driver"))}
	}
}

// Exec logs its params and calls the underlying init Exec method.
func (d *SlogDriver) Exec(ctx context.Context, query string, args, v any) error {
	d.log(ctx, "Exec", slog.String("query", query), slog.Any("args", args))
	return d.error(ctx, "Exec", d.Driver.Exec(ctx, query, args, v))
}

// ExecContext logs its params and calls the underlying init ExecContext method if it is supported.
func (d *SlogDriver) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	drv, ok := d.Driver.(interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.ExecContext is not supported")
	}
	d.log(ctx, "ExecContext", slog.String("query", query), slog.Any("args", args))
	result, err := drv.ExecContext(ctx, query, args...)
	return result, d.error(ctx, "ExecContext", err)
}

// Query logs its params and calls the underlying init Query method.
func (d *SlogDriver) Query(ctx context.Context, query string, args, v any) error {
	d.log(ctx, "Query", slog.String("query", query), slog.Any("args", args))
	return d.error(ctx, "Query", d.Driver.Query(ctx, query, args, v))
}

// QueryContext logs its params and calls the underlying init QueryContext method if it is supported.
func (d *SlogDriver) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	drv, ok := d.Driver.(interface {
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.QueryContext is not supported")
	}
	d.log(ctx, "QueryContext", slog.String("query", query), slog.Any("args", args))
	rows, err := drv.QueryContext(ctx, query, args...)
	return rows, d.error(ctx, "QueryContext", err)
}

// Tx adds an log-id for the transaction and calls the underlying init Tx command.
func (d *SlogDriver) Tx(ctx context.Context) (dialect.Tx, error) {
	tx, err := d.Driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	id := d.option.generateID(ctx)
	d.log(ctx, "Tx started", slog.String("id", id))
	return &SlogTx{Tx: tx, Handler: d.Handler.with(slog.String("database", "tx")), id: id, ctx: ctx}, nil
}

// BeginTx adds an log-id for the transaction and calls the underlying init BeginTx command if it is supported.
func (d *SlogDriver) BeginTx(ctx context.Context, opts *sql.TxOptions) (dialect.Tx, error) {
	drv, ok := d.Driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.BeginTx is not supported")
	}
	tx, err := drv.BeginTx(ctx, opts)
	if err != nil {
		return nil, d.error(ctx, "BeginTx", err)
	}
	id := d.option.generateID(ctx)
	d.log(ctx, "BeginTx started", slog.String("id", id))
	return &SlogTx{Tx: tx, Handler: d.Handler.with(slog.String("database", "tx")), id: id, ctx: ctx}, nil
}

// SlogTx is a transaction implementation that logs all transaction operations.
type SlogTx struct {
	dialect.Tx // underlying transaction.
	Handler
	id  string          // transaction logging id.
	ctx context.Context // underlying transaction context.
}

// Exec logs its params and calls the underlying transaction Exec method.
func (d *SlogTx) Exec(ctx context.Context, query string, args, v any) error {
	d.log(ctx, "Exec", slog.String("id", d.id), slog.String("query", query), slog.Any("args", args))
	return d.error(ctx, "Exec", d.Tx.Exec(ctx, query, args, v))
}

// ExecContext logs its params and calls the underlying transaction ExecContext method if it is supported.
func (d *SlogTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	drv, ok := d.Tx.(interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Tx.ExecContext is not supported")
	}
	d.log(ctx, "ExecContext", slog.String("id", d.id), slog.String("query", query), slog.Any("args", args))
	result, err := drv.ExecContext(ctx, query, args...)

	return result, d.error(ctx, "ExecContext", err)
}

// Query logs its params and calls the underlying transaction Query method.
func (d *SlogTx) Query(ctx context.Context, query string, args, v any) error {
	d.log(ctx, "Query", slog.String("id", d.id), slog.String("query", query), slog.Any("args", args))
	return d.error(ctx, "Query", d.Tx.Query(ctx, query, args, v))
}

// QueryContext logs its params and calls the underlying transaction QueryContext method if it is supported.
func (d *SlogTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	drv, ok := d.Tx.(interface {
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Tx.QueryContext is not supported")
	}
	d.log(ctx, "QueryContext", slog.String("id", d.id), slog.String("query", query), slog.Any("args", args))
	rows, err := drv.QueryContext(ctx, query, args...)

	return rows, d.error(ctx, "QueryContext", err)
}

// Commit logs this step and calls the underlying transaction Commit method.
func (d *SlogTx) Commit() error {
	d.log(d.ctx, "Commit", slog.String("id", d.id))
	return d.error(d.ctx, "Commit", d.Tx.Commit())
}

// Rollback logs this step and calls the underlying transaction Rollback method.
func (d *SlogTx) Rollback() error {
	d.log(d.ctx, "Rollback", slog.String("id", d.id))
	return d.error(d.ctx, "Rollback", d.Tx.Rollback())
}
