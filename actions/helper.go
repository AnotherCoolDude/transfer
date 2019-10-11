package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type query map[string]string

func responseBytes(response *http.Response) ([]byte, error) {
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()
	return bytes, nil
}

func unmarshal(response *http.Response, model interface{}) error {
	b, e := responseBytes(response)
	if e != nil {
		return e
	}
	e = json.Unmarshal(b, &model)
	if e != nil {
		return e
	}
	return nil
}
