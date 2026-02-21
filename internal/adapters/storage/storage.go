package storage

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/Util787/url-shortener/internal/common"
	"github.com/Util787/url-shortener/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ( // чтобы не загромождать конфиг
	defaultMaxConns        = 10
	defaultConnMaxLifetime = time.Hour
	defaultConnMaxIdleTime = time.Minute * 10
)

type PostgresStorage struct {
	pgxPool *pgxpool.Pool
}

func MustInitPostgres(ctx context.Context, cfg config.PostgresConfig) PostgresStorage {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)

	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(fmt.Errorf("failed to parse postgres connection string: %w", err))
	}

	// Pool configuration
	pgxConfig.MaxConns = defaultMaxConns
	pgxConfig.MaxConnLifetime = defaultConnMaxLifetime
	pgxConfig.MaxConnIdleTime = defaultConnMaxIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create postgres connection pool: %w", err))
	}

	err = pool.Ping(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to ping postgres: %w", err))
	}

	return PostgresStorage{
		pgxPool: pool,
	}
}

func (p *PostgresStorage) Shutdown() {
	p.pgxPool.Close()
}

func (p *PostgresStorage) SaveURL(ctx context.Context, id string, longURL string, shortURL string) error {
	op := common.GetOperationName()

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	upd := qb.Insert("url_mappings").
		Columns("long_url", "short_url", "created_at", "id").
		Values(longURL, shortURL, time.Now().Unix(), id)

	sqlStr, args, err := upd.ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build insert query: %w", op, err)
	}

	_, err = p.pgxPool.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to execute insert query: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) LongURLExists(ctx context.Context, longURL string) (bool, error) {
	op := common.GetOperationName()

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sel := qb.Select("short_url").
		From("url_mappings").
		Where(sq.Eq{"long_url": longURL})

	sqlStr, args, err := sel.ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: failed to build select query: %w", op, err)
	}

	var shortURL string
	err = p.pgxPool.QueryRow(ctx, sqlStr, args...).Scan(&shortURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("%s: failed to execute select query: %w", op, err)
	}

	return true, nil
}

func (p *PostgresStorage) ShortURLExists(ctx context.Context, shortURL string) (bool, error) {
	op := common.GetOperationName()

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sel := qb.Select("long_url").
		From("url_mappings").
		Where(sq.Eq{"short_url": shortURL})

	sqlStr, args, err := sel.ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: failed to build select query: %w", op, err)
	}

	var longURL string
	err = p.pgxPool.QueryRow(ctx, sqlStr, args...).Scan(&longURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("%s: failed to execute select query: %w", op, err)
	}

	return true, nil
}

func (p *PostgresStorage) GetRandomURL(ctx context.Context) (string, error) {
	op := common.GetOperationName()

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sel := qb.Select("short_url").
		From("url_mappings").
		OrderBy("RANDOM()").
		Limit(1)

	sqlStr, args, err := sel.ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: failed to build select query: %w", op, err)
	}

	var shortURL string
	err = p.pgxPool.QueryRow(ctx, sqlStr, args...).Scan(&shortURL)
	if err != nil {
		return "", fmt.Errorf("%s: failed to execute select query: %w", op, err)
	}

	return shortURL, nil
}

func (p *PostgresStorage) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	op := common.GetOperationName()

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sel := qb.Select("long_url").
		From("url_mappings").
		Where(sq.Eq{"short_url": shortURL})

	sqlStr, args, err := sel.ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: failed to build select query: %w", op, err)
	}

	var longURL string
	err = p.pgxPool.QueryRow(ctx, sqlStr, args...).Scan(&longURL)
	if err != nil {
		return "", fmt.Errorf("%s: failed to execute select query: %w", op, err)
	}

	return longURL, nil
}

func (p *PostgresStorage) DeleteURL(ctx context.Context, id *string, longURL *string, shortURL *string) error {
	op := common.GetOperationName()

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	del := qb.Delete("url_mappings")

	if id != nil {
		del = del.Where(sq.Eq{"id": *id})
	}
	if longURL != nil {
		del = del.Where(sq.Eq{"long_url": *longURL})
	}
	if shortURL != nil {
		del = del.Where(sq.Eq{"short_url": *shortURL})
	}

	sqlStr, args, err := del.ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build delete query: %w", op, err)
	}

	_, err = p.pgxPool.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to execute delete query: %w", op, err)
	}

	return nil
}
