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

type TrackRepository interface {
	GetTrackByID(ctx context.Context, id uuid.UUID) (*model.Track, error)
	GetTracksWithFilter(ctx context.Context, filter Filter) ([]model.Track, error)
	CreateTrack(ctx context.Context, track *model.Track) (*model.Track, error)
	CreateTracks(ctx context.Context, tracks []*model.Track) ([]*model.Track, error)
	UpdateTrack(ctx context.Context, track *model.Track) error
	ParitalUpdateTrack(ctx context.Context, track *model.Track) error
	DeleteTrack(ctx context.Context, id uuid.UUID) error
	DeleteTracks(ctx context.Context, id_list []uuid.UUID) error
	GetArtistOfTrack(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
}

type PostgresTrackRepository struct {
	dbpool *pgxpool.Pool
}

// UpdatePassword implements UserRepository.
func (tr *PostgresTrackRepository) UpdatePassword(ctx context.Context, id uuid.UUID, new_password string, old_password string) error {
	panic("unimplemented")
}

// UserLogin implements UserRepository.
func (tr *PostgresTrackRepository) UserLogin(ctx context.Context, email string, password string) error {
	panic("unimplemented")
}

// CreateUser implements UserRepository.
func (tr *PostgresTrackRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	panic("unimplemented")
}

// DeleteUser implements UserRepository.
func (tr *PostgresTrackRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// GetFollowArtist implements UserRepository.
func (tr *PostgresTrackRepository) GetFollowArtist(ctx context.Context, id uuid.UUID) ([]model.Artist, error) {
	panic("unimplemented")
}

// GetUserByID implements UserRepository.
func (tr *PostgresTrackRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	panic("unimplemented")
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

func (tr *PostgresTrackRepository) GetTracksWithFilter(ctx context.Context, filter Filter) ([]model.Track, error) {

	sort_criteria := filter.GetSortCriteria()
	fetchString := fmt.Sprintf(`
		select id, name, length from tracks
		where (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		order by %s id ASC
		limit $2 offset $3
	`, sort_criteria)

	rows, err := tr.dbpool.Query(ctx, fetchString, filter.Props["name"], filter.Limit, filter.GetOffSet())
	if err != nil {
		return nil, err
	}

	type OnlyTrackInfo struct {
		ID     uuid.UUID
		Name   string
		Length int
	}
	onlytracks, err := pgx.CollectRows(rows, pgx.RowToStructByName[OnlyTrackInfo])
	if err != nil {
		return nil, err
	}

	tracks := []model.Track{}
	for _, ot := range onlytracks {
		artists_id, err := tr.GetArtistOfTrack(ctx, ot.ID)
		if err != nil {
			return nil, err
		}
		track := model.Track{
			ID:       ot.ID,
			Name:     ot.Name,
			Length:   ot.Length,
			ArtistID: artists_id,
		}
		tracks = append(tracks, track)
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

	insert_xref_string := "INSERT INTO artists_tracks(artist_id, track_id) VALUES ($1, $2)"
	for _, artist_id := range track.ArtistID {
		_, err = tx.Exec(ctx, insert_xref_string, artist_id, track.ID)
		if err != nil {
			return nil, err
		}
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
	args := []any{
		track.ID,
		track.Name,
		track.Length,
	}
	_, err := tr.dbpool.Exec(ctx, "update tracks set name = $2, length = $3 where id = $1", args...)
	if err != nil {
		return err
	}

	return nil
}

func (tr *PostgresTrackRepository) ParitalUpdateTrack(ctx context.Context, track *model.Track) error {
	return nil
}

func (tr *PostgresTrackRepository) DeleteTrack(ctx context.Context, id uuid.UUID) error {
	deleteString := `
		with res as (DELETE FROM tracks where id = $1 returning 1)
		select count(*) from res
	`

	var success int
	if err := tr.dbpool.QueryRow(ctx, deleteString, id).Scan(&success); err != nil {
		return err
	}

	if success == 0 {
		return fmt.Errorf("there is no track with id: %s", id.String())
	}

	return nil
}

func (tr *PostgresTrackRepository) DeleteTracks(ctx context.Context, id_list []uuid.UUID) error {
	tx, err := tr.dbpool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	for _, id := range id_list {
		batch.Queue("delete tracks where id = $1", id)
	}

	br := tx.SendBatch(ctx, batch)

	br.Close()

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (tr *PostgresTrackRepository) GetArtistOfTrack(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	fetchString := "select artist_id from artists_tracks where track_id = $1"
	rows, err := tr.dbpool.Query(ctx, fetchString, id)
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
	return id_list, nil
}
