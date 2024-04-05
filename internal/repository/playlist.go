package repository

import (
	"context"
	"flotify/internal/model"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistRepository interface {
	CreatePlaylist(ctx context.Context, playlist model.Playlist) (*model.Playlist, error)
	GetPlaylist(ctx context.Context, playlist_id uuid.UUID) (*model.Playlist, error)
	AddTracksToPlaylist(ctx context.Context, playlist_id uuid.UUID, track_id_list []uuid.UUID) (*model.Playlist, error)
	DeleteTracksFromPlaylist(ctx context.Context, playlist_id uuid.UUID, track_id_list []uuid.UUID) (*model.Playlist, error)
	GetTracksOfPlaylist(ctx context.Context, playlist_id uuid.UUID) ([]model.Track, error)
	DeletePlaylist(ctx context.Context) error
}

type PostgresPlaylistRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresPlaylistRepository(dbpool *pgxpool.Pool) *PostgresPlaylistRepository {
	return &PostgresPlaylistRepository{
		dbpool: dbpool,
	}
}

func (pr *PostgresPlaylistRepository) CreatePlaylist(ctx context.Context, playlist model.Playlist) (*model.Playlist, error) {
	// let playlistname in db = playlist.name + user_id
	insertString := "INSERT INTO playlists(name, user_id) VALUES($1, $2) RETURNING id"

	var uuid_byte []byte
	err := pr.dbpool.QueryRow(ctx, insertString, playlist.UserID, playlist.Name).Scan(&uuid_byte)
	if err != nil {
		return nil, err
	}

	id, err := uuid.FromBytes(uuid_byte)
	if err != nil {
		return nil, err
	}
	playlist.ID = id

	return &playlist, nil
}

func (pr *PostgresPlaylistRepository) GetPlaylist(ctx context.Context, playlist_id uuid.UUID) (*model.Playlist, error) {
	tx, err := pr.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	playlist := model.Playlist{
		ID: playlist_id,
	}

	queryString := "SELECT name, user_id FROM playlists WHERE playlist_id=$1"
	err = tx.QueryRow(ctx, queryString, playlist_id).Scan(&playlist.Name, &playlist.UserID)
	if err != nil {
		return nil, err
	}

	queryString = "SELECT track_id FROM playlist_tracks WHERE playlist_id=$1"
	rows, err := tx.Query(ctx, queryString, playlist_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var track_id uuid.UUID
		err = rows.Scan(&track_id)
		if err != nil {
			return nil, err
		}
		playlist.TrackIDList = append(playlist.TrackIDList, track_id)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &playlist, nil
}

func (pr *PostgresPlaylistRepository) AddTracksToPlaylist(ctx context.Context, playlist_id uuid.UUID, track_id_list []uuid.UUID) (*model.Playlist, error) {
	tx, err := pr.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	context, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	type track_playlist struct {
		playlist_id uuid.UUID
		track_id    uuid.UUID
	}
	rows := [][]any{}
	for _, track_id := range track_id_list {
		rows[0] = append(rows[0], track_playlist{playlist_id: playlist_id, track_id: track_id})
	}

	_, err = tx.CopyFrom(
		context,
		pgx.Identifier{"playlists_tracks"},
		[]string{"playlist_id", "track_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return nil, err
	}

	playlist, err := pr.GetPlaylist(ctx, playlist_id)
	if err != nil {
		return nil, err
	}

	playlist.TrackIDList = append(playlist.TrackIDList, track_id_list...)
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func (pr *PostgresPlaylistRepository) DeleteTracksFromPlaylist(ctx context.Context, playlist_id uuid.UUID, track_id_list []uuid.UUID) (*model.Playlist, error) {
	tx, err := pr.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	context, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	type track_playlist struct {
		playlist_id uuid.UUID
		track_id    uuid.UUID
	}

	batch := &pgx.Batch{}

	for _, track_id := range track_id_list {
		batch.Queue("delete tracks where track_id = $1", track_id)
	}

	br := tx.SendBatch(context, batch)

	br.Close()

	playlist, err := pr.GetPlaylist(ctx, playlist_id)
	if err != nil {
		return nil, err
	}

	delete_track_id := make(map[uuid.UUID]bool, len(track_id_list))
	for _, track_id := range track_id_list {
		delete_track_id[track_id] = true
	}

	var new_track_id_list []uuid.UUID
	for _, track_id := range playlist.TrackIDList {
		if !delete_track_id[track_id] {
			new_track_id_list = append(new_track_id_list, track_id)
		}
	}
	playlist.TrackIDList = new_track_id_list

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func (pr *PostgresPlaylistRepository) GetTracksOfPlaylist(ctx context.Context, playlist_id uuid.UUID) ([]model.Track, error) {
	return []model.Track{}, nil
}

func (pr *PostgresPlaylistRepository) DeletePlaylist(ctx context.Context) error {
	return nil
}
