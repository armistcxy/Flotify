package helper

import (
	"errors"
	"flotify/internal/config"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

type AuthManager struct {
	SecretKey string
}

type AuthCredential struct {
	ID uuid.UUID
}

func NewAuthManager() AuthManager {
	auth_config := config.LoadAuthConfig()
	return AuthManager{
		auth_config.SecretKey,
	}
}
func (am *AuthManager) GenerateJWT(ac *AuthCredential, exp_time time.Duration) (string, error) {

	expiration_time := time.Now().UTC().Add(exp_time)

	jwtclaim := jwt.MapClaims{
		"id":  ac.ID,
		"exp": expiration_time.Unix(),
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwtclaim,
	)

	if token_string, err := token.SignedString([]byte(am.SecretKey)); err != nil {
		return "", err
	} else {
		return token_string, nil
	}
}

func (am *AuthManager) VerifyJWT(token_string string, id uuid.UUID) error {
	log.Println(token_string)
	token, err := jwt.Parse(
		token_string,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(am.SecretKey), nil
		},
	)
	// there something that I unexpected here: when parsing, the exp is auto validated
	// but I want to handle it manual
	if err != nil {
		log.Println("Not parse")
		return err
	}
	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			exp := int64(claims["exp"].(float64))
			if exp < time.Now().UTC().Unix() {
				return errors.New("token expired")
			}

			id_string_form := claims["id"].(string)
			id_check, err := uuid.FromString(id_string_form)
			if err != nil {
				return err
			}

			if id_check != id {
				return errors.New("wrong id")
			}
			return nil
		} else {
			return errors.New("can't retrieve claims from token")
		}

	} else {
		return errors.New("invalid token")
	}
}
