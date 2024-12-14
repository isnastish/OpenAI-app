package db

import "context"

// NOTE: We want to have an interface in DB with Postgres being one of them
// The interface should support main methods for writing and reading the data
// So we could swith between SQL and NoSQL data storages.
// But for now lets stick with a single implementation.
// Later, when we have a clear structure, we could replace Postgres
// with MySQL for example, or keep both.

type DatabaseController interface {
	AddUser(ctx context.Context) error
	HasUser(ctx context.Context) (bool, error)
	Close() error
}
