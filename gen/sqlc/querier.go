// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package sqlc

import (
	"context"
)

type Querier interface {
	DeleteUser(ctx context.Context, id int32) error
	InsertUser(ctx context.Context, arg InsertUserParams) (int32, error)
	SelectUsers(ctx context.Context) ([]SelectUsersRow, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
}

var _ Querier = (*Queries)(nil)
