package sshelter

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ncruces/zenity"
	"github.com/oxodao/sshelter_client/config"
	"github.com/oxodao/sshelter_client/models"
)

type Client struct {
	Info       func(interface{})
	baseUrl    *url.URL
	httpClient *http.Client
	cfg        *config.Config
}

func New(cfg *config.Config, log func(interface{})) (*Client, error) {
	u, err := url.Parse(cfg.Server.Url)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return &Client{
		baseUrl:    u,
		cfg:        cfg,
		Info:       log,
		httpClient: http.DefaultClient,
	}, err
}

func buildUserAgent() string {
	return "sshelter/" + config.VERSION + " - " + runtime.GOOS + "/" + runtime.GOARCH
}

func (c *Client) newRequest(method, path string, body interface{}, authenticated bool) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseUrl.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	c.Info(method + " @ " + u.String())
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/ld+json")
	req.Header.Set("User-Agent", buildUserAgent())

	if authenticated {
		req.Header.Set("Authorization", "Bearer "+c.cfg.Server.Token)
	}

	return req, nil
}

func (c *Client) authenticatedRequest(method, path string, body interface{}) (*http.Request, error) {
	return c.newRequest(method, path, body, true)
}

func (c *Client) unauthenticatedRequest(method, path string, body interface{}) (*http.Request, error) {
	return c.newRequest(method, path, body, false)
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)

	if resp != nil && resp.StatusCode >= 400 {
		var msg map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&msg)
		if err != nil {
			return resp, err
		}

		fmt.Println(msg)
		return resp, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

//#region Authentication
func (c *Client) shouldAuthenticate() bool {
	if len(c.cfg.Server.Token) == 0 {
		return true
	}

	splitted := strings.Split(c.cfg.Server.Token, ".")
	if len(splitted) != 3 {
		return true
	}

	token := splitted[1]
	tokenDecoded, err := base64.StdEncoding.DecodeString(token + "==")
	if err != nil {
		return true
	}

	var tokenData map[string]interface{}
	err = json.Unmarshal(tokenDecoded, &tokenData)
	if err != nil {
		return true
	}

	exp := tokenData["exp"].(float64) - 30 // if the token is expiring in less than 30 seconds

	return float64(time.Now().Unix()) > exp
}

func (c *Client) AuthenticateIfRequired() error {
	if !c.shouldAuthenticate() {
		return nil
	}

	if len(c.cfg.Server.RefreshToken) != 0 {
		c.Info("Refreshing token")

		req, err := c.unauthenticatedRequest("POST", "api/auth/refresh", map[string]string{
			"refresh_token": c.cfg.Server.RefreshToken,
		})
		if err != nil {
			return err
		}

		var msg map[string]interface{}
		res, err := c.do(req, &msg)
		if err != nil {
			// hahahahahahahahahahahahahahahahaha
			// no but seriously, handle this properly later

			// If the server is not joignable, the user should NOT be asked for his password,
			// but the software will be used in "offline mode" (The user can want to port forward on his LAN for example)

			goto authentication
		}

		if res.StatusCode != http.StatusOK {
			fmt.Println(res.StatusCode)
			goto authentication
		}

		c.cfg.Server.Token = msg["token"].(string)
		c.cfg.Server.RefreshToken = msg["refresh_token"].(string)

		return c.cfg.Save()
	}

authentication:
	// Ask the user for his password & log him in
	c.Info("Authenticating")

	loggedIn := false

	for !loggedIn {
		usr, pwd, err := zenity.Password(zenity.Title("Login to sshelter"), zenity.Username())
		if err != nil {
			if err == zenity.ErrCanceled {
				zenity.Error("Authentication cancelled", zenity.Title("Login to sshelter"))
				os.Exit(0)
			} else {
				panic(err)
			}
		}

		req, _ := c.unauthenticatedRequest("POST", "api/auth/login", map[string]string{
			"username": usr,
			"password": pwd,
		})

		var data map[string]interface{}
		res, err := c.do(req, &data)
		if err != nil {
			// Same as in the refresh,
			// If the server is not joignable, the user should NOT be asked for his password,
			// but the software will be used in "offline mode" (The user can want to port forward on his LAN for example)

			return err
		}

		if res.StatusCode == http.StatusUnauthorized {
			zenity.Error("Wrong username or password", zenity.Title("Login to sshelter"))
			continue
		}

		if res.StatusCode != http.StatusOK {
			// @TODO set offline mode flag in PRV without cyclic depedencies (:
			return fmt.Errorf("Unexpected status code: %d", res.StatusCode)
		}

		// Could be useful later to pre-fill the username field
		c.cfg.Server.Username = usr
		c.cfg.Server.Token = data["token"].(string)
		c.cfg.Server.RefreshToken = data["refresh_token"].(string)

		loggedIn = true
	}

	return c.cfg.Save()
}

//#endregion

func (c *Client) GetMachines() ([]models.Machine, error) {
	err := c.AuthenticateIfRequired()
	if err != nil {
		return nil, err
	}

	req, err := c.authenticatedRequest("GET", "api/machines", nil)
	if err != nil {
		return nil, err
	}

	var msg map[string]json.RawMessage
	_, err = c.do(req, &msg)
	if err != nil {
		return nil, err
	}

	machines := []models.Machine{}
	err = json.Unmarshal(msg["hydra:member"], &machines)

	return machines, err
}

// half of the following methods are generated through copilot *VISAGE CHOQUÉ*
func (c *Client) CreateMachine(m *models.Machine) error {
	err := c.AuthenticateIfRequired()
	if err != nil {
		return err
	}

	if m.ForwardedPorts == nil {
		m.ForwardedPorts = []models.ForwardedPort{}
	}

	if m.Port == 0 {
		m.Port = 22
	}

	req, err := c.authenticatedRequest("POST", "api/machines", m)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteMachine(machine *models.Machine) error {
	err := c.AuthenticateIfRequired()
	if err != nil {
		return err
	}

	req, err := c.authenticatedRequest("DELETE", *machine.Id, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
