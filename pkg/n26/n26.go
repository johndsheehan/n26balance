package n26

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const n26URL = "https://api.tech26.de"

// Balance n26 balance details
type Balance struct {
	AvailableBalance float64 `json:"availableBalance"`
	UsableBalance    float64 `json:"usableBalance"`
	IBAN             string  `json:"iban"`
	BIC              string  `json:"bic"`
	BankName         string  `json:"bankName"`
	Seized           bool    `json:"seized"`
	ID               string  `json:"id"`
}

// Config n26 username and password
type Config struct {
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
}

// N26 client connection
type N26 struct {
	client *http.Client
}

// NewClient new n26 client connection
func NewClient(config Config) (N26, error) {
	n26Client := N26{}

	c := oauth2.Config{
		ClientID:     "android",
		ClientSecret: "secret",
		Endpoint: oauth2.Endpoint{
			TokenURL: n26URL + "/oauth/token",
		},
	}

	ctx := context.Background()
	tkn, err := c.PasswordCredentialsToken(ctx, config.Username, config.Password)
	if err != nil {
		return n26Client, err
	}

	n26Client = N26{
		client: c.Client(ctx, tkn),
	}
	return n26Client, nil
}

// Balance fetch n26 balance
func (n *N26) Balance() (Balance, error) {
	balance := Balance{}

	u, _ := url.ParseRequestURI(n26URL)
	u.Path = "/api/accounts"

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return balance, err
	}

	rsp, err := n.client.Do(req)
	if err != nil {
		return balance, err
	}
	defer rsp.Body.Close()

	bdy, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return balance, err
	}

	err = json.Unmarshal(bdy, &balance)
	if err != nil {
		return balance, err
	}

	return balance, nil
}
