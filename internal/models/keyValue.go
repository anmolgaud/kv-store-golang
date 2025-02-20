package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type CacheEntry struct {
	ID int64 `json:"id"`
	Key string `json:"key"`
	Value string `json:"value"`
	TTL int64 `json:"ttl"`
	IsDeleted int `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
}

type KeyValueModel struct {
	DB *sql.DB
}

func (m KeyValueModel) Insert(cache *CacheEntry) error {
	nowMilli := time.Now().UTC().UnixMilli()
	query := `INSERT INTO kv (key, value, ttl) VALUES (?, ?, ?);`
	args := []any{
		cache.Key,
		cache.Value,
		nowMilli + cache.TTL,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err;
}

func (m KeyValueModel) Get(key string) (*CacheEntry, error) {
	query := `SELECT id, key, value, ttl, is_deleted, created_at FROM kv WHERE ttl >= ? AND key = ?;`
	cacheEntry := CacheEntry{}

	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second);
	defer cancel()
	nowMilli := time.Now().UTC().UnixMilli()
	err := m.DB.QueryRowContext(ctx, query, nowMilli, key).Scan(
		&cacheEntry.ID,
		&cacheEntry.Key,
		&cacheEntry.Value,
		&cacheEntry.TTL,
		&cacheEntry.IsDeleted,
		&cacheEntry.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &cacheEntry, nil;
}

func (m KeyValueModel) Delete(key string) error {
	query := `UPDATE kv SET is_deleted = 1 WHERE key = ?;`
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second);
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, key);
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}