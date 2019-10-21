package actions

import (
	"crypto/tls"
	"github.com/AnotherCoolDude/transfer/models"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
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

func (c *proadclient) fetchProject(projectno string, project *models.PAProject) error {
	projectresp, err := c.do("GET", "projects", http.NoBody, query{"projectno": projectno})
	if err != nil {
		return err
	}
	var pp []models.PAProject
	err = unmarshalProad(projectresp, &pp)
	*project = pp[0]
	if err != nil {
		return err
	}
	return nil
}

func (c *proadclient) fetchProjectAsync(projectno string, project *models.PAProject, sem chan int, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	sem <- 1
	if err := c.fetchProject(projectno, project); err != nil {
		select {
		case errChan <- err:
			// we are the first worker to fail
		default:
			// there allready happend an error
		}
	}
	<-sem
}

// func (c *proadclient) fetchProjects(projectnumbers []string) ([]models.PAProject, error) {
// 	projects := []models.PAProject{}
// 	for _, projectno := range projectnumbers {
// 		projectresp, err := c.do("GET", "projects", http.NoBody, query{"projectno": projectno})
// 		if err != nil {
// 			return projects, err
// 		}
// 		var project models.PAProject
// 		err = unmarshalProad(projectresp, &project)
// 		if err != nil {
// 			return projects, err
// 		}
// 		projects = append(projects, project)
// 	}
// 	return projects, nil
// }

func (c *proadclient) fetchTodos(project *models.PAProject) error {
	todosresp, err := c.do("GET", "tasks", http.NoBody, query{"project": strconv.Itoa(project.Urno)})
	if err != nil {
		return err
	}
	var todos []models.PATodo
	err = unmarshalProad(todosresp, &todos)
	if err != nil {
		return err
	}
	project.Todos = todos
	return nil
}

func (c *proadclient) fetchTodosAsync(project *models.PAProject, sem chan int, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	sem <- 1
	if err := c.fetchTodos(project); err != nil {
		select {
		case errChan <- err:
			// we are the first worker to fail
		default:
			// there allready happend an error
		}
	}
	<-sem
}

// func (c *proadclient) asyncUnmarshal(URL string, query map[string]string, unmarshalChan chan asyncUnmarshal) {
// 	resp, err := c.do("GET", URL, http.NoBody, query)

// 	if err != nil {
// 		unmarshalChan <- asyncUnmarshal{err: err}
// 		return
// 	}
// 	au := asyncUnmarshal{}
// 	au.err = unmarshalProad(resp, &au.model)
// 	for _, v := range query {
// 		au.breadcrumb = v
// 	}
// 	unmarshalChan <- au
// }

// type asyncResponse struct {
// 	response *http.Response
// 	err      error
// }

// type asyncUnmarshal struct {
// 	model      interface{}
// 	err        error
// 	breadcrumb string
// }

// func (au *asyncUnmarshal) decode(model interface{}) error {
// 	if au.err != nil {
// 		return au.err
// 	}
// 	return mapstructure.Decode(au.model, &model)
// }
