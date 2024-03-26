package repository

import (
	"context"
	"flotify/internal/model"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUserInfo(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetFollowArtist(ctx context.Context, id uuid.UUID) ([]model.Artist, error)
}

type PostgresUserRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresUserRepository(dbpool *pgxpool.Pool) *PostgresTrackRepository {
	return &PostgresTrackRepository{
		dbpool: dbpool,
	}
}

func (ur *PostgresUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := ur.dbpool.QueryRow(ctx, "select name from tracks where id=$1", id)

	user := model.User{ID: id}
	err := row.Scan(&user.Username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *PostgresUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	tx, err := ur.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	context, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	hash_password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return nil, err
	}

	insertString := "INSERT INTO users(name) VALUES ($1) RETURNING id"
	args := []any{
		user.Username,
		hash_password,
	}
	row := tx.QueryRow(context, insertString, args...)

	uuid_byte := []byte{}
	if err = row.Scan(&uuid_byte); err != nil {
		return nil, err
	}

	user.ID, err = uuid.FromBytes(uuid_byte)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (tr *PostgresTrackRepository) UpdateUserInfo(ctx context.Context, user *model.User) error {
	args := []any{
		user.ID,
		user.Username,
	}

	_, err := tr.dbpool.Exec(ctx, "update users set name = $2 where id = $1", args...)
	if err != nil {
		return err
	}

	return nil
}

func (ur *PostgresUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	deleteString := `
		with res as (DELETE FROM users where id = $1 returning 1)
		select count(*) from res
	`

	var success int
	if err := ur.dbpool.QueryRow(ctx, deleteString, id).Scan(&success); err != nil {
		return err
	}

	if success == 0 {
		return fmt.Errorf("there is no user with id: %s", id.String())
	}

	return nil
}

func (ur *PostgresUserRepository) GetFollowArtist(ctx context.Context, id uuid.UUID) ([]model.Artist, error) {
	fetchString := `
		SELECT (*) FROM artists
		WHERE id BELONG TO (SELECT artist_id FROM artists_users WHERE user_id = $1)
	`

	rows, err := ur.dbpool.Query(ctx, fetchString, id)
	if err != nil {
		return nil, err
	}

	artists := []model.Artist{}
	for rows.Next() {
		artist := model.Artist{}
		err = rows.Scan(&artist.ID, &artist.Name, &artist.Description)
		if err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}

	return artists, nil

}
