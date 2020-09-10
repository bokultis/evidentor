package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/bokultis/evidentor/api/apiutil"
	"github.com/bokultis/evidentor/api/authorization"
	"github.com/bokultis/evidentor/api/db"
	"github.com/bokultis/evidentor/api/logger"
	"github.com/bokultis/evidentor/api/token"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gorilla/context"
)

// JWTMiddleware for auth
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//init casbin auth rules
		a, err := gormadapter.NewAdapterByDB(db.DB)
		if err != nil {
			log.Fatal(err)
		}
		authEnforcer, err := casbin.NewEnforcer("auth_model.conf", a)
		if err != nil {
			log.Fatal(err)
		}

		//claims, err := VerifyToken(r)
		tokenAuth, err := token.ExtractTokenMetadata(r)
		if err != nil {
			logger.Logger.Warn(err.Error())
			apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
			return
		}

		userID, err := token.FetchAuth(tokenAuth)
		if err != nil {
			logger.Logger.Warn(err.Error())
			apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
			return
		}
		context.Set(r, "user_id", userID)
		authorization.Authorizer(authEnforcer, strconv.FormatUint(userID, 10))(next).ServeHTTP(w, r)
	})
}
