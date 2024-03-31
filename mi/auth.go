// MiAuth のドキュメントは[Misskey Hub]と[Misskey Forum]を参照。
//
// [Misskey Hub]: https://misskey-hub.net/ja/docs/for-developers/api/token/miauth/
// [Misskey Forum]: https://forum.misskey.io/d/6-miauth
package mi

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

// Callback URLにはリダイレクト時に session クエリパラメータにセッションIDが付加されます。
//
// Permission で利用可能な権限の一覧は[こちら]
//
// [こちら]: https://misskey-hub.net/ja/docs/for-developers/api/permission/
type AuthConfig struct {
	SessionID  uuid.UUID // SessionID はUUID v4 推奨。省略すると自動生成します。
	Name       string    // Name はクライアントアプリケーション名を指定します。
	Callback   string    // Callback は認証後のリダイレクト先URLを指定します。
	Permission []string  // Permission は使用したい権限を指定します。
	Host       string    // Host はサーバのURLをスキーマを含めて指定します。
}

type AuthResp struct {
	OK    bool   // OK は認証の成否です。
	Token string // Token です。
	User  User   // User にはユーザ情報が返ります。
}

// AuthCodeURL は認証ページのURLを返します。
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

// Exchange は Token とユーザ情報を含む AuthResp を返します。
//
// ユーザによる認証後に実行します。
func (c *AuthConfig) Exchange(ctx context.Context) (AuthResp, error) {
	u, err := url.Parse(c.Host)
	if err != nil {
		panic(err)
	}

	u = u.JoinPath("api", "miauth", c.SessionID.String(), "check")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return AuthResp{}, err
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return AuthResp{}, err
	}
	defer resp.Body.Close()

	var r AuthResp
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return AuthResp{}, err
	}
	// pp.Println(r)

	return r, nil
}
