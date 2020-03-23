package user

import (
	"database/sql"
	"encoding/json"
	"evidentor/api/apiutil"
	"evidentor/api/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func UsersIndexHandler(w http.ResponseWriter, r *http.Request) {

	var users []*UserDO
	users, err := GetAllUsers()
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}
	logger.Logger.Warn("test")
	usersWO := make([]interface{}, len(users))
	for i, user := range users {
		usersWO[i] = makeUserWO(user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersWO)
}

func UsersShowHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])
	user, err := GetUser(userId)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(makeUserWO(user))
}

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
		token, err := loggingUserWO.generateJWT()
		if err != nil {
			logger.Logger.Warn(err.Error())
			apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
			return
		}
		json.NewEncoder(w).Encode(&token)
	} else {
		apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
		return
	}
}

func UsersDeleteHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	userId, _ := strconv.Atoi(params["userId"])
	err := DeleteUser(userId)

	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("User deleted")
}

func UsersUpdateHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var user UserUpdateWO
	userId, _ := strconv.Atoi(params["userId"])

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

	err = UpdateUser(userId, *update)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}

	updatedUser, err := GetUser(userId)
	if err != nil {
		logger.Logger.Warn(err.Error())
		apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(makeUserWO(updatedUser))
}
