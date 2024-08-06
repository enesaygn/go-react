package repository

import (
	"sasa-elterminali-service/internal/adapters/cache"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	postgres *pgxpool.Pool
	cache    *cache.RedisCache
}

// new database
func NewDB(postgres *pgxpool.Pool, cache *cache.RedisCache) *DB {
	return &DB{
		postgres: postgres,
		cache:    cache,
	}
}
