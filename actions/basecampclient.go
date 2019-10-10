package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AnotherCoolDude/transfer/models"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

type basecampclient struct {
	email       string
	appName     string
	id          int
	oauthConfig *oauth2.Config
	state       string
	code        string
	token       *oauth2.Token
	httpclient  *http.Client
}

func defaultBasecampclient() *basecampclient {
	return &basecampclient{
		appName: os.Getenv("BASECAMP_APPNAME"),
		email:   os.Getenv("BASECAMP_EMAIL"),
		id:      0,
		state:   xid.New().String(),
		code:    "",
		token:   &oauth2.Token{},
		oauthConfig: &oauth2.Config{
			RedirectURL:  os.Getenv("BASECAMP_CALLBACK"),
			ClientID:     os.Getenv("BASECAMP_CLIENT"),
			ClientSecret: os.Getenv("BASECAMP_SECRET"),
			Scopes:       []string{},
			Endpoint: oauth2.Endpoint{
				AuthStyle: oauth2.AuthStyleAutoDetect,
				AuthURL:   "https://launchpad.37signals.com/authorization/new",
				TokenURL:  "https://launchpad.37signals.com/authorization/token",
			},
		},
		httpclient: http.DefaultClient,
	}
}

func (c *basecampclient) do(method, URL string, body io.Reader) (*http.Response, error) {
	requestURL, err := url.Parse(URL)
	if err != nil {
		fmt.Println("could not parse url")
		return nil, err
	}
	if !requestURL.IsAbs() {
		requestURL = c.baseURL()
		requestURL.Path = path.Join(requestURL.Path, URL)
	}
	req, err := http.NewRequest(method, requestURL.String(), body)
	if err != nil {
		fmt.Println("error making request " + err.Error())
		return nil, err
	}
	c.addHeader(req)
	resp, err := c.httpclient.Do(req)
	if err != nil {
		fmt.Println("error making request " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (c *basecampclient) baseURL() *url.URL {
	urlString := fmt.Sprintf("https://3.basecampapi.com/%d/", c.id)
	url, _ := url.Parse(urlString)
	return url
}

func (c *basecampclient) addHeader(request *http.Request) {
	request.Header.Add("Authorization", "Bearer "+c.token.AccessToken)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", fmt.Sprintf("%s (%s)", c.appName, c.email))
}

// AuthCodeURL returns the url where the client needs to authenticate
func (c *basecampclient) AuthCodeURL() string {
	return c.oauthConfig.AuthCodeURL(c.state, oauth2.SetAuthURLParam("type", "web_server"))
}

func (c *basecampclient) handleCallback(request *http.Request) error {
	code := request.FormValue("code")
	state := request.FormValue("state")
	if state != c.state {
		return errors.New("[basecamp.go/handleCallback] state doesn't match")
	}
	t, err := c.oauthConfig.Exchange(oauth2.NoContext, code, oauth2.SetAuthURLParam("type", "web_server"))
	if err != nil {
		return err
	}
	c.token = t
	return nil
}

// IsValid returns wether the the basecamp client has a valid token
func (c *basecampclient) IsValid() bool {
	if !c.token.Valid() {
		return false
	}
	if c.id == 0 {
		return false
	}
	return true
}

func (c *basecampclient) receiveID() error {
	if !c.token.Valid() {
		return errors.New("need valid token to receive ID")
	}
	resp, err := c.do("GET", "https://launchpad.37signals.com/authorization.json", http.NoBody)
	if err != nil {
		return errors.New("[basecamp.go/receiveID] couldn't make request to auth endpoint")
	}
	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("[basecamp.go/receiveID] couldn't read response body")
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.Unmarshal(respbytes, &result)
	accounts := result["accounts"].([]interface{})
	accDetails := accounts[0].(map[string]interface{})
	c.id = int(accDetails["id"].(float64))
	return nil
}

func unmarshalBCProjects(response *http.Response) ([]models.BCProject, error) {
	var bcprojects []models.BCProject
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return bcprojects, err
	}
	defer response.Body.Close()

	err = json.Unmarshal(bytes, &bcprojects)
	if err != nil {
		return bcprojects, err
	}

	return bcprojects, nil
}
