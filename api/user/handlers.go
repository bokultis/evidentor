package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bokultis/evidentor/api/apiutil"
	"github.com/bokultis/evidentor/api/logger"
	"github.com/bokultis/evidentor/api/token"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//UsersIndexHandler index handler
func UsersIndexHandler(w http.ResponseWriter, r *http.Request) {

	var users []*UserDO
	users, err := GetAllUsers()
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	usersWO := make([]interface{}, len(users))
	for i, user := range users {
		usersWO[i] = makeUserWO(user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersWO)
}

//UsersShowHandler show user by ID
func UsersShowHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	userID, _ := strconv.Atoi(params["userId"])
	user, err := GetUser(userID)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(makeUserWO(user))
}

//UsersCreateHandler create new user
func UsersCreateHandler(w http.ResponseWriter, r *http.Request) {

	var user UserInputWO
	var firstName, lastName, address sql.NullString
	var birthday mysql.NullTime

	err := apiutil.GetRequestBody(r, &user)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
	}

	validationErr := user.validate()
	if validationErr != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, validationErr)
		return
	}

	if user.FirstName == nil {
		firstName = sql.NullString{Valid: false, String: ""}
	} else {
		firstName = sql.NullString{Valid: true, String: *user.FirstName}
	}

	if user.LastName == nil {
		lastName = sql.NullString{Valid: false, String: ""}
	} else {
		lastName = sql.NullString{Valid: true, String: *user.LastName}
	}

	if user.Address == nil {
		address = sql.NullString{Valid: false, String: ""}
	} else {
		address = sql.NullString{Valid: true, String: *user.Address}
	}

	if user.Birthday == nil {
		birthday = mysql.NullTime{Valid: false, Time: time.Now()}
	} else {
		bt, _ := time.Parse("2006-01-02", *user.Birthday)
		birthday = mysql.NullTime{Valid: true, Time: bt}
	}

	password := hashPassword(*user.Password)

	newUser := &UserDO{
		FirstName: firstName,
		LastName:  lastName,
		Gender:    *user.Gender,
		Birthday:  birthday,
		Email:     *user.Email,
		Address:   address,
		Optin:     *user.Optin,
		Password:  password,
	}

	userID, err := CreateUser(newUser)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	userCreated, err := GetUser(userID)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	//fmt.Printf("%+v ", userCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(makeUserWO(userCreated))
}

//UsersDeleteHandler delete user
func UsersDeleteHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	userID, _ := strconv.Atoi(params["userId"])
	err := DeleteUser(userID)

	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("User deleted")
}

//UsersUpdateHandler update user
func UsersUpdateHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var user UserUpdateWO
	userID, _ := strconv.Atoi(params["userId"])

	err := apiutil.GetRequestBody(r, &user)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
	}

	update, validationErr := user.validate()
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, validationErr)
		return
	}

	err = UpdateUser(userID, *update)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	updatedUser, err := GetUser(userID)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(makeUserWO(updatedUser))
}

// LoginHandler user login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var uCredentials UserCredentials

	err := apiutil.GetRequestBody(r, &uCredentials)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	if uCredentials.Email == nil || uCredentials.Password == nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
		return
	}

	loggingUserDO, err := GetUserByEmail(*uCredentials.Email)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
		return
	}

	loggingUserWO := makeUserWO(loggingUserDO)

	w.Header().Set("Content-Type", "application/json")

	if loggingUserDO.checkPassword(*uCredentials.Password) {
		generatedToken, err := loggingUserWO.generateJWT()
		if err != nil {
			logger.Logger.Warn(err.Error())
			apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
			return
		}
		err = createAuth(loggingUserWO.ID, &generatedToken)
		if err != nil {
			logger.Logger.Warn(err.Error())
			apiutil.NewErrorResponse(w, apiutil.ErrBadParameter)
			return
		}

		tokens := map[string]string{
			"access_token":  generatedToken.AccessToken,
			"refresh_token": generatedToken.RefreshToken,
		}

		json.NewEncoder(w).Encode(&tokens)
	} else {
		apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
		return
	}
}

// LogoutHandler user logout handler
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	au, err := token.ExtractTokenMetadata(r)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
		return
	}
	deleted, delErr := token.DeleteAuth(au.AccessUUID)
	if delErr != nil || deleted == 0 { //if any goes wrong
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Successfully logged out")
}

//RefreshTokenHandler token func
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	mapToken := map[string]string{}

	err := apiutil.GetRequestBody(r, &mapToken)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}
	refreshToken := mapToken["refresh_token"]
	logger.Logger.Warn(refreshToken)

	//verify the token
	signingKey := []byte(os.Getenv("JWT_SECRET"))

	rToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
		return
	}
	//is token valid?
	if _, ok := rToken.Claims.(jwt.Claims); !ok && !rToken.Valid {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := rToken.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && rToken.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			//c.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userID, err := strconv.Atoi(fmt.Sprintf("%.f", claims["user_id"]))
		if err != nil {
			logger.Logger.Warn(err.Error())
			apiutil.NewErrorResponse(w, apiutil.ErrBadParameter)
			return
		}

		//Delete the previous Refresh Token
		deleted, delErr := token.DeleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			logger.Logger.Warn(err.Error())
			apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
			return
		}
		loggingUserDO, err := GetUser(userID)
		if err != nil {
			logger.Logger.Warn(err.Error())
			apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
			return
		}

		loggingUserWO := makeUserWO(loggingUserDO)
		//Create new pairs of refresh and access tokens
		ts, createErr := loggingUserWO.generateJWT()
		if createErr != nil {
			//c.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		saveErr := createAuth(userID, &ts)
		if saveErr != nil {
			//c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		json.NewEncoder(w).Encode(&tokens)
	} else {
		apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
		return
	}
}
