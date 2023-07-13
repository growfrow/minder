// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: signing_keys.sql

package db

import (
	"context"
	"time"
)

const createSigningKey = `-- name: CreateSigningKey :one
INSERT INTO signing_keys (group_id, private_key, public_key, passphrase, key_identifier, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, group_id, private_key, public_key, passphrase, key_identifier, created_at, updated_at
`

type CreateSigningKeyParams struct {
	GroupID       int32     `json:"group_id"`
	PrivateKey    string    `json:"private_key"`
	PublicKey     string    `json:"public_key"`
	Passphrase    string    `json:"passphrase"`
	KeyIdentifier string    `json:"key_identifier"`
	CreatedAt     time.Time `json:"created_at"`
}

func (q *Queries) CreateSigningKey(ctx context.Context, arg CreateSigningKeyParams) (SigningKey, error) {
	row := q.db.QueryRowContext(ctx, createSigningKey,
		arg.GroupID,
		arg.PrivateKey,
		arg.PublicKey,
		arg.Passphrase,
		arg.KeyIdentifier,
		arg.CreatedAt,
	)
	var i SigningKey
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.PrivateKey,
		&i.PublicKey,
		&i.Passphrase,
		&i.KeyIdentifier,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteSigningKey = `-- name: DeleteSigningKey :exec
DELETE FROM signing_keys WHERE group_id = $1 AND key_identifier = $2
`

type DeleteSigningKeyParams struct {
	GroupID       int32  `json:"group_id"`
	KeyIdentifier string `json:"key_identifier"`
}

func (q *Queries) DeleteSigningKey(ctx context.Context, arg DeleteSigningKeyParams) error {
	_, err := q.db.ExecContext(ctx, deleteSigningKey, arg.GroupID, arg.KeyIdentifier)
	return err
}

const getSigningKeyByGroupID = `-- name: GetSigningKeyByGroupID :one
SELECT id, group_id, private_key, public_key, passphrase, key_identifier, created_at, updated_at FROM signing_keys WHERE group_id = $1
`

func (q *Queries) GetSigningKeyByGroupID(ctx context.Context, groupID int32) (SigningKey, error) {
	row := q.db.QueryRowContext(ctx, getSigningKeyByGroupID, groupID)
	var i SigningKey
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.PrivateKey,
		&i.PublicKey,
		&i.Passphrase,
		&i.KeyIdentifier,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSigningKeyByIdentifier = `-- name: GetSigningKeyByIdentifier :one
SELECT id, group_id, private_key, public_key, passphrase, key_identifier, created_at, updated_at FROM signing_keys WHERE key_identifier = $1
`

func (q *Queries) GetSigningKeyByIdentifier(ctx context.Context, keyIdentifier string) (SigningKey, error) {
	row := q.db.QueryRowContext(ctx, getSigningKeyByIdentifier, keyIdentifier)
	var i SigningKey
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.PrivateKey,
		&i.PublicKey,
		&i.Passphrase,
		&i.KeyIdentifier,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
