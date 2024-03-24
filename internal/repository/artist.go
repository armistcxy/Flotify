package repository

import (
	"context"
	"flotify/internal/model"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtistRepository interface {
	GetArtistByID(ctx context.Context, id uuid.UUID) (*model.Artist, error)
	GetArtistsWithFilter(ctx context.Context, filter Filter) ([]model.Artist, error)
	CreateArtist(ctx context.Context, artist *model.Artist) (*model.Artist, error)
	CreateArtists(ctx context.Context, artists []*model.Artist) ([]*model.Artist, error)
	UpdateArtist(ctx context.Context, artist *model.Artist) error
	PartialUpdateArtist(ctx context.Context, artist *model.Artist) error
	DeleteArtist(ctx context.Context, id uuid.UUID) error
	DeleteArtists(ctx context.Context, id_list []uuid.UUID) error
	GetTrackOfArtist(ctx context.Context, id uuid.UUID) ([]*model.Track, error)
}

type PostgresArtistRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresArtistRepository(dbpool *pgxpool.Pool) *PostgresArtistRepository {
	return &PostgresArtistRepository{
		dbpool: dbpool,
	}
}

func (ar *PostgresArtistRepository) GetArtistByID(ctx context.Context, id uuid.UUID) (*model.Artist, error) {
	row := ar.dbpool.QueryRow(ctx, "select id, name, description from artists where id=$1", id)

	artist := model.Artist{}

	uuid_byte := []byte{}
	err := row.Scan(&uuid_byte, &artist.Name, &artist.Description)
	if err != nil {
		return nil, err
	}

	id, err = uuid.FromBytes(uuid_byte)
	if err != nil {
		return nil, err
	}
	artist.ID = id
	return &artist, nil
}

func (ar *PostgresArtistRepository) GetArtistsWithFilter(ctx context.Context, filter Filter) ([]model.Artist, error) {

	sort_criteria := filter.GetSortCriteria()
	fetchString := fmt.Sprintf(`
        select id, name, description from artists
        where (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
        order by %s id ASC
		limit $2 offset $3
    `, sort_criteria)

	rows, err := ar.dbpool.Query(ctx, fetchString, filter.Props["name"], filter.Limit, filter.GetOffSet())
	if err != nil {
		return nil, err
	}

	artists, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Artist])
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (ar *PostgresArtistRepository) CreateArtist(ctx context.Context, artist *model.Artist) (*model.Artist, error) {
	tx, err := ar.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	context, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	insertString := "INSERT INTO artists(name, description) VALUES ($1, $2) RETURNING id"
	args := []any{
		artist.Name,
		artist.Description,
	}
	row := tx.QueryRow(context, insertString, args...)

	uuid_byte := []byte{}
	if err = row.Scan(&uuid_byte); err != nil {
		return nil, err
	}

	artist.ID, err = uuid.FromBytes(uuid_byte)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return artist, err
}

func (ar *PostgresArtistRepository) CreateArtists(ctx context.Context, artists []*model.Artist) ([]*model.Artist, error) {
	tx, err := ar.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	context, cancel := context.WithTimeout(ctx, time.Millisecond*20)
	defer cancel()

	// use bulk insert through COPY protocol

	rows := [][]any{}
	for _, artist := range artists {
		rows[0] = append(rows[0], artist)
	}
	_, err = tx.CopyFrom(
		context,
		pgx.Identifier{"artists"},
		[]string{"name", "description"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (ar *PostgresArtistRepository) UpdateArtist(ctx context.Context, artist *model.Artist) error {
	args := []any{
		artist.ID,
		artist.Name,
		artist.Description,
	}
	_, err := ar.dbpool.Exec(ctx, "update artists set name = $2, description = $3 where id = $1", args...)
	if err != nil {
		return err
	}

	return nil
}

func (ar *PostgresArtistRepository) PartialUpdateArtist(ctx context.Context, artist *model.Artist) error {
	return nil
}

func (ar *PostgresArtistRepository) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	deleteString := `
        with res as (DELETE FROM artists where id = $1 returning 1)
        select count(*) from res
    `

	var success int
	if err := ar.dbpool.QueryRow(ctx, deleteString, id).Scan(&success); err != nil {
		return err
	}

	if success == 0 {
		return fmt.Errorf("there is no artist with id: %s", id.String())
	}

	return nil
}

func (ar *PostgresArtistRepository) DeleteArtists(ctx context.Context, id_list []uuid.UUID) error {
	tx, err := ar.dbpool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	for _, id := range id_list {
		batch.Queue("delete artists where id = $1", id)
	}

	br := tx.SendBatch(ctx, batch)

	br.Close()

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (ar *PostgresArtistRepository) GetTrackOfArtist(ctx context.Context, id uuid.UUID) ([]*model.Track, error) {
	fetchString := "select track_id from artists_tracks where artist_id = $1"
	rows, err := ar.dbpool.Query(ctx, fetchString, id)
	if err != nil {
		return nil, err
	}

	id_list := []uuid.UUID{}
	for rows.Next() {
		uuid_byte := []byte{}
		err = rows.Scan(&uuid_byte)
		if err != nil {
			return nil, err
		}
		id, err := uuid.FromBytes(uuid_byte)
		if err != nil {
			return nil, err
		}
		id_list = append(id_list, id)
	}

	track := []*model.Track{}
	return track, nil
}
