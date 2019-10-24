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
	for i := range todos {
		todos[i].Project = project
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

func (c *proadclient) CreateTodoFromBasecamp(basecampTodo models.BCTodo, proadProject models.PAProject) error {
	type postTask struct {
			UrnoManager      int    `json:"urno_manager"`
			UrnoCompany      int    `json:"urno_company"`
			UrnoProject      int    `json:"urno_project"`
			UrnoServiceCode  int    `json:"urno_service_code"`
			UrnoResponsible  int    `json:"urno_responsible"`
			Shortinfo        string `json:"shortinfo"`
			FromDatetime     string `json:"from_datetime"`
			UntilDatetime    string `json:"until_datetime"`
			ReminderDatetime string `json:"reminder_datetime"`
			Status           string `json:"status"`
			Priority         string `json:"priority"`
			Description      string `json:"description"`
	}
	
	fromdate := basecampTodo.CreatedAt.Format("2006-01-02T15:04:05")
	pt := postTask{UrnoCompany: proadProject.}
}
