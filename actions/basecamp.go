package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/rs/xid"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

type client struct {
	email       string
	appName     string
	id          int
	oauthConfig *oauth2.Config
	state       string
	code        string
	token       *oauth2.Token
	httpclient  *http.Client
}

func defaultClient() *client {
	return &client{
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

func (c *client) do(method, URL string, body io.Reader) (*http.Response, error) {
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

func (c *client) baseURL() *url.URL {
	urlString := fmt.Sprintf("https://3.basecampapi.com/%d/", c.id)
	url, _ := url.Parse(urlString)
	return url
}

func (c *client) addHeader(request *http.Request) {
	request.Header.Add("Authorization", "Bearer "+c.token.AccessToken)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", fmt.Sprintf("%s (%s)", c.appName, c.email))
}

func (c *client) authCodeURL() string {
	return c.oauthConfig.AuthCodeURL(c.state, oauth2.SetAuthURLParam("type", "web_server"))
}

func (c *client) handleCallback(request *http.Request) {
	code := request.FormValue("code")
	state := request.FormValue("state")
	if state != c.state {
		fmt.Println("[basecamp.go/handleCallback] state doesn't match")
		return
	}
	t, err := c.oauthConfig.Exchange(oauth2.NoContext, code, oauth2.SetAuthURLParam("type", "web_server"))
	if err != nil {
		fmt.Println("[basecamp.go/handleCallback] couldn't exchange token: " + err.Error())
		return
	}
	c.token = t
}

func (c *client) receiveID() error {
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

var (
	bcClient = defaultClient()
)

// BasecampShow default implementation.
func BasecampShow(c buffalo.Context) error {
	authURL := bcClient.authCodeURL()
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// BasecampCallback handles the callback from basecamp upon authentication
func BasecampCallback(c buffalo.Context) error {
	bcClient.handleCallback(c.Request())

	if err := bcClient.receiveID(); err != nil {
		return c.Render(404, r.HTML("index.html"))
	}

	resp, err := bcClient.do("GET", "/projects.json", http.NoBody)
	if err != nil {
		fmt.Println(err)
		return c.Render(404, r.HTML("index.html"))
	}
	defer resp.Body.Close()
	bb, _ := ioutil.ReadAll(resp.Body)
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, bb, "", "\t")
	if err != nil {
		fmt.Println(err)
		return c.Render(404, r.HTML("index.html"))
	}
	return c.Render(200, r.Func("application/json", func(w io.Writer, d render.Data) error {
		_, err := w.Write(bb)
		return err
	}))

	// return c.Render(200, r.String(prettyJSON.String()))
}
