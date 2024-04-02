package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flotify/internal/custom_error"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

type AuthManager struct {
	SecretKey  string
	repository AuthRepository
}

type AuthCredential struct {
	ID uuid.UUID
}

func NewAuthManager(secretkey string, repository AuthRepository) AuthManager {
	return AuthManager{
		SecretKey:  secretkey,
		repository: repository,
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
				return custom_error.AccessTokenExpiredError{}
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
		return custom_error.InvalidTokenError{}
	}
}

func (am *AuthManager) StoreRefreshToken(user_id uuid.UUID, refresh_token string) error {
	err := am.repository.AddRefreshToken(user_id, refresh_token)
	if err != nil {
		return err
	}

	return nil
}

// EncryptRefreshToken encrypts the refresh token using AES encryption
func EncryptRefreshToken(refreshToken string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(refreshToken))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(refreshToken))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptRefreshToken decrypts the encrypted refresh token using AES encryption
func DecryptRefreshToken(encryptedToken string, key []byte) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return string(ciphertext), nil
}
