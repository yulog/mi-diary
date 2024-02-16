package mi

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type AuthConfig struct {
	SessionID  uuid.UUID
	Name       string
	Callback   string
	Permission []string
	Host       string // スキーマを含める
}

type AuthResp struct {
	OK    bool
	Token string
	User  User
}

func (c *AuthConfig) AuthCodeURL() string {
	u, err := url.Parse(c.Host)
	if err != nil {
		panic(err)
	}

	if c.SessionID == uuid.Nil {
		c.SessionID, err = uuid.NewRandom()
		if err != nil {
			panic(err)
		}
	}

	u = u.JoinPath("miauth", c.SessionID.String())

	q, _ := url.ParseQuery("")
	q.Set("name", c.Name)
	q.Set("permission", strings.Join(c.Permission, ","))
	if c.Callback != "" {
		q.Set("callback", c.Callback)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func (c *AuthConfig) Exchange() (AuthResp, error) {
	u, err := url.Parse(c.Host)
	if err != nil {
		panic(err)
	}

	u = u.JoinPath("api", "miauth", c.SessionID.String(), "check")

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return AuthResp{}, err
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return AuthResp{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AuthResp{}, err
	}

	var r AuthResp
	json.Unmarshal(body, &r)
	// pp.Println(r)

	return r, nil
}
