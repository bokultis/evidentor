package user

import (
	"regexp"
	"evidentor/api/apiutil"
	"time"

	"github.com/go-sql-driver/mysql"
)

func (usr *UserInputWO) validate() *apiutil.Error {
	var errItems []*apiutil.ErrorItem

	rEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString

	if usr.Password == nil {
		errItems = append(errItems, &apiutil.ErrorItem{Name: "password",
			Message: "{Password must be provided}",
		})
	}

	if usr.Email == nil {
		errItems = append(errItems, &apiutil.ErrorItem{Name: "email",
			Message: "{Email must be provided}",
		})
	} else {
		matched := rEmail(*usr.Email)
		if !matched {
			errItems = append(errItems, &apiutil.ErrorItem{Name: "email",
				Message: "Email wrong format",
			})
		}

		user, _ := GetUserByEmail(*usr.Email)

		if user != nil {
			errItems = append(errItems, &apiutil.ErrorItem{Name: "email",
				Message: "Email is taken",
			})
		}

	}

	if usr.Gender == nil {
		gender := "unknown"
		usr.Gender = &gender
	} else {
		switch *usr.Gender {
		case "male", "female", "other", "unknown":
		default:
			errItems = append(errItems, &apiutil.ErrorItem{Name: "gender",
				Message: "Gender must be in (male,female,other,unknown)",
			})
		}
	}

	if usr.Birthday == nil {
	} else {
		brd := *usr.Birthday // + "T00:00:01.000Z"
		_, err := time.Parse("2006-01-02", brd)
		if err != nil {
			errItems = append(errItems, &apiutil.ErrorItem{Name: "birthday",
				Message: "Birthday date invalid 1" + err.Error(),
			})
		}

	}

	if usr.Optin == nil {
		optin := "unknown"
		usr.Optin = &optin
	} else {
		switch *usr.Optin {
		case "unknown", "yes", "no":
		default:
			errItems = append(errItems, &apiutil.ErrorItem{Name: "optin",
				Message: "Optin must be in (unknown,yes,no)",
			})
		}
	}

	if len(errItems) > 0 {
		return apiutil.NewValidationError(errItems)
	}

	return nil
}

func (usr *UserUpdateWO) validate() (*map[string]interface{}, *apiutil.Error) {

	rEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString

	var errItems []*apiutil.ErrorItem
	update := map[string]interface{}{}

	if usr.FirstName != nil {
		update["first_name"] = *usr.FirstName
	}

	if usr.LastName != nil {
		update["last_name"] = *usr.LastName
	}

	if usr.Gender != nil {
		switch *usr.Gender {
		case "male", "female", "other", "unknown":
			update["gender"] = *usr.Gender
		default:
			errItems = append(errItems, &apiutil.ErrorItem{Name: "gender",
				Message: "invalid value"})
		}
	}

	if usr.Birthday != nil { // to do
		if *usr.Birthday == "" {
			update["birthday"] = mysql.NullTime{Valid: false, Time: time.Now()}
		} else {
			update["birthday"] = *usr.Birthday
		}
	}

	if usr.Email != nil { //
		if *usr.Email == "" {
			update["e_mail"] = *usr.Email
		} else {
			matched := rEmail(*usr.Email)
			if !matched {
				errItems = append(errItems, &apiutil.ErrorItem{Name: "email",
					Message: "invalid value"})
			} else {
				update["e_mail"] = *usr.Email
			}
		}

	}

	if usr.Address != nil { //

		update["address"] = *usr.Address

	}

	if usr.Optin != nil {
		switch *usr.Optin {
		case "unknown", "yes", "no":
			update["optin"] = *usr.Optin
		default:
			errItems = append(errItems, &apiutil.ErrorItem{Name: "optin",
				Message: "invalid value"})
		}
	}

	if len(errItems) > 0 {
		return nil, apiutil.NewValidationError(errItems)
	}

	return &update, nil
}
