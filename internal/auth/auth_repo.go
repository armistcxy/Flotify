package auth

import (
	"context"
	"flotify/internal/custom_error"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	dbpool    *pgxpool.Pool
	secretkey string
}

func NewAuthRepository(dbpool *pgxpool.Pool, secretkey string) *AuthRepository {
	return &AuthRepository{
		dbpool:    dbpool,
		secretkey: secretkey,
	}
}

func (ar *AuthRepository) VerifyRefreshToken(user_id, refresh_token string) error {
	token, err := jwt.Parse(
		refresh_token,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(ar.secretkey), nil
		},
	)
	if err != nil {
		return err
	}

	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			exp := int64(claims["exp"].(float64))
			if exp < time.Now().UTC().Unix() {
				return custom_error.RefreshTokenExpired{}
			}
		}
	}

	queryString := "SELECT token FROM refreshtokens WHERE user_id=$1"
	var db_token string
	err = ar.dbpool.QueryRow(context.Background(), queryString, user_id).Scan(&db_token)
	if err != nil {
		return err
	}

	// db_token, err := DecryptRefreshToken(encrypt_db_token, []byte(ar.secretkey))
	// if err != nil {
	// 	return err
	// }

	if refresh_token != db_token {
		return custom_error.InvalidTokenError{}
	}

	return nil
}

func (ar *AuthRepository) AddRefreshToken(user_id uuid.UUID, refresh_token string) error {
	// encrypt_refresh_token, err := EncryptRefreshToken(refresh_token, []byte(ar.secretkey))
	// if err != nil {
	// 	return err
	// }

	insertString := "INSERT INTO refreshtokens (user_id, token) VALUES($1, $2)"
	_, err := ar.dbpool.Exec(context.Background(), insertString, user_id, refresh_token)
	if err != nil {
		return err
	}
	return nil
}
