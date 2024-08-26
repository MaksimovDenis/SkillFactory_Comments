package storage

import (
	"database/sql"

	"context"

	"git.mobiledep.ru/flagshtok/backend/controller/storage/queries"
	_ "github.com/jackc/pgx/v4/stdlib" // pgx
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// Storage postgresql datastore wrapper
type Storage struct {
	dbh     *sql.DB
	l       zerolog.Logger
	pg      *pgxpool.Pool
	Queries *queries.Queries
}

// NewStorage create and init new storage instance
func NewStorage(ctx context.Context, pgConnString string, log zerolog.Logger) (*Storage, error) {
	pgConn, err := pgxpool.New(ctx, pgConnString)
	if err != nil {
		log.Error().Err(err).Msg("failed to init postgres connection")
		return nil, err
	}

	if err = pgConn.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("failed to connect to postgres db")
		pgConn.Close()
		return nil, err
	}

	if err := migration(pgConnString); err != nil {
		return nil, err
	}

	dbConn, err := sql.Open("pgx", pgConnString)
	if err != nil {
		return nil, err
	}

	dbConn.SetMaxOpenConns(5)
	dbConn.SetMaxIdleConns(1)

	hdl := &Storage{
		dbh: dbConn,
		l:   log,
		pg:  pgConn,
	}

	return db, nil
}

func (repo *Repository) StopPG() {
	if repo.PostgresDB != nil {
		repo.Log.Info().Msg("closing PostgreSQL connection pool")
		repo.PostgresDB.Close()
	}
}
