package middleware

import (
	"evidentor/api/apiutil"
	"evidentor/api/authorization"
	"evidentor/api/db"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gorilla/context"

	jwt "github.com/dgrijalva/jwt-go"
)

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

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := VerifyToken(tokenString)
		if err != nil {
			apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
			return
		}
		//fmt.Printf("%+v", claims.(jwt.MapClaims)["user_id"].(float64))
		//pass userId claim to req
		//todo: find a better way to convert the claim to string
		userId := strconv.FormatFloat(claims.(jwt.MapClaims)["user_id"].(float64), 'g', 1, 64)
		context.Set(r, "user_id", userId)
		authorization.Authorizer(authEnforcer, userId)(next).ServeHTTP(w, r)
	})
}
