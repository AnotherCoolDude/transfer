package actions

import (
	"crypto/tls"
	"encoding/json"
	"github.com/AnotherCoolDude/transfer/models"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

type proadclient struct {
	httpClient *http.Client
	apiKey     string
}

func defaultProadclient() *proadclient {
	return &proadclient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		apiKey: os.Getenv("PROAD_APIKEY"),
	}
}

func (c *proadclient) do(method, URL string, body io.Reader, query map[string]string) (*http.Response, error) {
	requestURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	if !requestURL.IsAbs() {
		requestURL, _ = url.Parse("https://192.168.0.15/api/v5/")
		requestURL.Path = path.Join(requestURL.Path, URL)
	}
	req, err := http.NewRequest(method, requestURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("apikey", c.apiKey)
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *proadclient) async(method, URL string, body io.Reader, query map[string]string, responseChanel chan asyncResponse) {
	resp, err := c.do(method, URL, body, query)
	responseChanel <- asyncResponse{response: resp, err: err}
}

func unmarshalPAProjects(response *http.Response) ([]models.PAProject, error) {
	type projectlist struct {
		Projects []models.PAProject `json:"project_list"`
	}
	var pl projectlist

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return pl.Projects, err
	}
	defer response.Body.Close()

	err = json.Unmarshal(bytes, &pl)
	if err != nil {
		return pl.Projects, err
	}

	return pl.Projects, nil
}

func unmarshalPATodos(response *http.Response) ([]models.PATodo, error) {
	type todolist struct {
		Todos []models.PATodo `json:"todo_list"`
	}
	var tl todolist

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return tl.Todos, err
	}
	defer response.Body.Close()

	err = json.Unmarshal(bytes, &tl)
	if err != nil {
		return tl.Todos, err
	}

	return tl.Todos, nil
}

type asyncResponse struct {
	response *http.Response
	err      error
}
