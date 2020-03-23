package apiutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//GetRequestBody decodes JSON in body to obj
func GetRequestBody(r *http.Request, obj interface{}) error {
	respBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBodyBytes, obj)
	if err != nil {
		return err
	}
	return nil
}
