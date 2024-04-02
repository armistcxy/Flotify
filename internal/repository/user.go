package repository

import (
	"context"
	"flotify/internal/custom_error"
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
	UpdatePassword(ctx context.Context, id uuid.UUID, new_password, old_password string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetFollowArtist(ctx context.Context, id uuid.UUID) ([]model.Artist, error)
	UserLogin(ctx context.Context, email string, password string) (*uuid.UUID, error)
}

type PostgresUserRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresUserRepository(dbpool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		dbpool: dbpool,
	}
}

func (ur *PostgresUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := ur.dbpool.QueryRow(ctx, "select username, email from users where id=$1", id)

	user := model.User{ID: id}
	err := row.Scan(&user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *PostgresUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	check_exist_string := "select count(id) from users where id = $1"

	var count int
	err := ur.dbpool.QueryRow(ctx, check_exist_string, user.ID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count != 0 {
		return nil, fmt.Errorf("username %s has already been used", user.Username)

	}

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

	insertString := "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING id"
	args := []any{
		user.Username,
		user.Email,
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

// this function need authentication
func (tr *PostgresUserRepository) UpdateUserInfo(ctx context.Context, user *model.User) error {
	args := []any{
		user.ID,
		user.Username,
	}

	check_exist_string := "select count(*) from users where name = $1"
	var count int

	if err := tr.dbpool.QueryRow(ctx, check_exist_string, user.Username).Scan(&count); err != nil {
		return err
	} else if count != 0 {
		return fmt.Errorf("username %s has already been used", user.Username)
	}
	_, err := tr.dbpool.Exec(ctx, "update users set name = $2 where id = $1", args...)
	if err != nil {
		return err
	}

	return nil
}

// this function need authentication, and verify using email
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

// this function need authentication
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

// this function need authentication
func (ur *PostgresUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, new_password, old_password string) error {
	// first get the hash password stored in db
	get_hash_password_string := "select password from users where id = $1"
	var hash_password string
	if err := ur.dbpool.QueryRow(ctx, get_hash_password_string, id).Scan(&hash_password); err != nil {
		return err
	}
	// next confirm old password
	if err := bcrypt.CompareHashAndPassword([]byte(hash_password), []byte(old_password)); err != nil {
		return custom_error.OldPasswordMismatchError{}
	}
	// then change passsword
	change_password_string := "update table users set password = $1 where id = $2"
	if _, err := ur.dbpool.Exec(ctx, change_password_string, new_password, id); err != nil {
		return err
	}
	return nil // nil mean change password success
}

func (ur *PostgresUserRepository) UserLogin(ctx context.Context, email string, password string) (*uuid.UUID, error) {
	// first determine whether the email in db
	check_exist_email_string := "select id from users where email = $1"
	var uuid_byte []byte
	if err := ur.dbpool.QueryRow(ctx, check_exist_email_string, email).Scan(&uuid_byte); err != nil {
		return nil, custom_error.MismatchError{}
	}
	// next get the hash password stored in db
	get_hash_password_string := "select password from users where email = $1"
	var hash_password string
	if err := ur.dbpool.QueryRow(ctx, get_hash_password_string, email).Scan(&hash_password); err != nil {
		return nil, err
	}
	// next check the password
	if err := bcrypt.CompareHashAndPassword([]byte(hash_password), []byte(password)); err != nil {
		return nil, custom_error.MismatchError{} // don't let the hacker know what was wrong
	}

	uuid, err := uuid.FromBytes(uuid_byte)
	if err != nil {
		return nil, err
	}
	return &uuid, nil
}
