package repository

import (
	"context"
	"flotify/internal/model"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistRepository interface {
	CreatePlaylist(ctx context.Context, playlist model.Playlist) (*model.Playlist, error)
	AddTracksToPlaylist(ctx context.Context, track_id_list []uuid.UUID) (*model.Playlist, error)
	DeleteTracksFromPlaylist(ctx context.Context, track_id_list []uuid.UUID) (*model.Playlist, error)
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

func (pr *PostgresPlaylistRepository) AddTracksToPlaylist(ctx context.Context, track_id_list []uuid.UUID) (*model.Playlist, error) {
	return nil, nil
}

func (pr *PostgresPlaylistRepository) DeleteTracksFromPlaylist(ctx context.Context, track_id_list []uuid.UUID) (*model.Playlist, error) {
	return nil, nil
}

func (pr *PostgresPlaylistRepository) GetTracksOfPlaylist(ctx context.Context, playlist_id uuid.UUID) ([]model.Track, error) {
	return []model.Track{}, nil
}

func (pr *PostgresPlaylistRepository) DeletePlaylist(ctx context.Context) error {
	return nil
}
