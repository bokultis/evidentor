package student

import (
	"evidentor/api/apiutil"
	"time"

	"github.com/go-sql-driver/mysql"
)

func (st *StudentInputWO) validate() *apiutil.Error {
	var errItems []*apiutil.ErrorItem

	if st.Gender == nil {
		gender := "unknown"
		st.Gender = &gender
	} else {
		switch *st.Gender {
		case "male", "female", "unknown":
		default:
			errItems = append(errItems, &apiutil.ErrorItem{Name: "gender",
				Message: "Gender must be in (male,female,unknown)",
			})
		}
	}

	if st.Birthday == nil {
	} else {
		brd := *st.Birthday // + "T00:00:01.000Z"
		_, err := time.Parse("2006-01-02", brd)
		if err != nil {
			errItems = append(errItems, &apiutil.ErrorItem{Name: "birthday",
				Message: "Birthday date invalid :" + err.Error(),
			})
		}

	}

	if len(errItems) > 0 {
		return apiutil.NewValidationError(errItems)
	}

	return nil
}

func (st *StudentUpdateWO) validate() (*map[string]interface{}, error) {
	var errItems []*apiutil.ErrorItem
	update := map[string]interface{}{}

	if st.FirstName != nil {
		update["first_name"] = *st.FirstName
	}

	if st.LastName != nil {
		update["last_name"] = *st.LastName
	}

	if st.Gender != nil {
		switch *st.Gender {
		case "male", "female", "unknown":
			update["gender"] = *st.Gender
		default:
			errItems = append(errItems, &apiutil.ErrorItem{Name: "gender",
				Message: "Gender must be in (male,female,unknown)",
			})
		}
	}

	if st.Birthday != nil { // to do
		if *st.Birthday == "" {
			update["birthday"] = mysql.NullTime{Valid: false, Time: time.Now()}
		} else {
			update["birthday"] = *st.Birthday
		}
	}

	if len(errItems) > 0 {
		return nil, apiutil.NewValidationError(errItems)
	}

	return &update, nil
}
