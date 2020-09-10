package authorization

import (
	"net/http"
	"strconv"

	"github.com/bokultis/evidentor/api/apiutil"
	"github.com/bokultis/evidentor/api/logger"
	"github.com/bokultis/evidentor/api/user"

	"github.com/casbin/casbin/v2"
)

func Authorizer(e *casbin.Enforcer, userId string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.NewLogger()
			role := "anonymous"
			userID, err := strconv.Atoi(userId)
			if err != nil {
				apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
				return
			}

			userRole, err := user.GetUserRole(userID)
			if err != nil {
				log.Warn(err.Error())
				apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
				return
			}

			if len(userRole.Name) > 0 {
				role = userRole.Name
			}
			//fmt.Println(userRole.Name)
			//fmt.Println(r.URL.Path)
			// casbin rule enforcing
			res, err := e.Enforce(role, r.URL.Path, r.Method)
			if err != nil {
				log.Printf("Enforcer return error %v", err)
				apiutil.NewErrorResponse(w, apiutil.NewIntError(err))
				return
			}
			if res {
				next.ServeHTTP(w, r)
			} else {
				log.Printf("ACCESS FORBIDEN")
				apiutil.NewErrorResponse(w, apiutil.ErrNotAuthenticated)
				return
			}
		}

		return http.HandlerFunc(fn)
	}
}
