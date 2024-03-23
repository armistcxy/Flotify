package repository

import (
	"context"
	"flotify/internal/model"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackRepository interface {
	GetTrackByID(ctx context.Context, id uuid.UUID) (*model.Track, error)
	GetAllTrack(ctx context.Context) ([]model.Track, error)
	CreateTrack(ctx context.Context, track *model.Track) (*model.Track, error)
	CreateTracks(ctx context.Context, tracks []*model.Track) ([]*model.Track, error)
	UpdateTrack(ctx context.Context, track *model.Track) error
	ParitalUpdateTrack(ctx context.Context, track *model.Track) error
	DeleteTrack(ctx context.Context, id uuid.UUID) error
	DeleteTracks(ctx context.Context, id_list []uuid.UUID) error
}

type PostgresTrackRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresTrackRepository(dbpool *pgxpool.Pool) *PostgresTrackRepository {
	return &PostgresTrackRepository{
		dbpool: dbpool,
	}
}

func (tr *PostgresTrackRepository) GetTrackByID(ctx context.Context, id uuid.UUID) (*model.Track, error) {
	row := tr.dbpool.QueryRow(ctx, "select id, name, length from tracks where id=$1", id)

	track := model.Track{}

	uuid_byte := []byte{}
	err := row.Scan(&uuid_byte, &track.Name, &track.Length)
	if err != nil {
		return nil, err
	}

	id, err = uuid.FromBytes(uuid_byte)
	if err != nil {
		return nil, err
	}
	track.ID = id
	return &track, nil
}

func (tr *PostgresTrackRepository) GetAllTrack(ctx context.Context) ([]model.Track, error) {
	rows, err := tr.dbpool.Query(ctx, "select id, name, length from track")
	if err != nil {
		return nil, err
	}

	tracks, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Track])
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (tr *PostgresTrackRepository) CreateTrack(ctx context.Context, track *model.Track) (*model.Track, error) {
	tx, err := tr.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	context, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	insertString := "INSERT INTO tracks(name, length) VALUES ($1, $2) RETURNING id"
	args := []any{
		track.Name,
		track.Length,
	}
	row := tx.QueryRow(context, insertString, args...)

	uuid_byte := []byte{}
	if err = row.Scan(&uuid_byte); err != nil {
		return nil, err
	}

	track.ID, err = uuid.FromBytes(uuid_byte)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return track, err
}

func (tr *PostgresTrackRepository) CreateTracks(ctx context.Context, tracks []*model.Track) ([]*model.Track, error) {
	tx, err := tr.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	context, cancel := context.WithTimeout(ctx, time.Millisecond*20)
	defer cancel()

	// use bulk insert throught COPY protocol

	rows := [][]any{}
	for _, track := range tracks {
		rows[0] = append(rows[0], track)
	}
	_, err = tx.CopyFrom(
		context,
		pgx.Identifier{"tracks"},
		[]string{"name", "length"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (tr *PostgresTrackRepository) UpdateTrack(ctx context.Context, track *model.Track) error {
	return nil
}

func (tr *PostgresTrackRepository) ParitalUpdateTrack(ctx context.Context, track *model.Track) error {
	return nil
}

func (tr *PostgresTrackRepository) DeleteTrack(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (tr *PostgresTrackRepository) DeleteTracks(ctx context.Context, id_list []uuid.UUID) error {
	return nil
}
