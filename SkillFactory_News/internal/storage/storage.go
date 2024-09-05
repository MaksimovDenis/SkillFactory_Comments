package storage

import (
	"context"

	"github.com/MaksimovDenis/SkillFactory_News/internal/models"
	feedsdb "github.com/MaksimovDenis/SkillFactory_News/internal/storage/feedsDB"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

type Feeds interface {
	StoreFeeds(feeds []models.Feeds) error                                              //Save news to storage
	Feeds(limit int, page int, title string) (feeds []models.Feeds, err error)          //Get news from storge
	FeedById(id int) (*models.Feeds, error)                                             //Get news by Id from storge
	FeedsByFilter(limit int, page int, filter string) (feeds []models.Feeds, err error) //Get news by filter and limit from storge
}

type Storage struct {
	Feeds      Feeds
	Log        zerolog.Logger
	PostgresDB *pgxpool.Pool
}

func NewRepository(ctx context.Context, postgresDB *pgxpool.Pool, log zerolog.Logger) *Storage {
	return &Storage{
		Feeds:      feedsdb.NewFeedsPostgres(postgresDB, log),
		Log:        log,
		PostgresDB: postgresDB,
	}
}
