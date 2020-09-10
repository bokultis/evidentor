package token

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bokultis/evidentor/api/logger"
	"github.com/bokultis/evidentor/api/redis"
	"github.com/dgrijalva/jwt-go"
)

//AccessDetails contain user asseccs data
type AccessDetails struct {
	AccessUUID string
	UserID     uint64
}

// VerifyToken verify JWT token based on environment var (JWT_SECRET)
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET"))
	tokenString := ExtractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token, err
}

// TokenValid validetes token
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

//ExtractToken retrive token from header
func ExtractToken(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		return ""
	}
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	return tokenString
}

//ExtractTokenMetadata extract metadata from token
func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		logger.Logger.Warn(err.Error())
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}

//FetchAuth get redis record
func FetchAuth(authD *AccessDetails) (uint64, error) {
	userid, err := redis.RedisClient.Get(redis.Ctx, authD.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

//DeleteAuth from redis
func DeleteAuth(givenUUID string) (int64, error) {
	deleted, err := redis.RedisClient.Del(redis.Ctx, givenUUID).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
