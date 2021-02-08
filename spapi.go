package spapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

type Engine struct {
	config Config
}

type Config struct {
	Region  string
	SandBox bool
	Beta    bool
	Credentials
}

type Credentials struct {
	ClientID          string
	SPAPIClientID     string
	SPAPIClientSecret string
	SPAPICallbackURL  string
	AWSAccessKeyID    string
	AWSAccessKey      string
	AccessToken       string
	RefreshToken      string
	TokenType         string
	ExpiresIn         int64
}

type OAuthSPAPI struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

func New(config Config) (*Engine, error) {

	engine := Engine{}

	err := checkRegion(config.Region)
	if err != nil {
		return nil, err
	}

	if config.Credentials.ClientID == "" {
		return nil, errors.New("Credentials.ClientID cannot be empty")
	}

	if config.Credentials.SPAPIClientID == "" {
		return nil, errors.New("Credentials.SPAPIClientID cannot be empty")
	}

	if config.Credentials.SPAPIClientSecret == "" {
		return nil, errors.New("Credentials.SPAPIClientSecret cannot be empty")
	}

	if config.Credentials.SPAPICallbackURL == "" {
		return nil, errors.New("Credentials.SPAPICallbackURL cannot be empty")
	}

	if config.Credentials.AWSAccessKeyID == "" {
		return nil, errors.New("Credentials.AWSAccessKeyID cannot be empty")
	}

	if config.Credentials.AWSAccessKey == "" {
		return nil, errors.New("Credentials.AWSAccessKey cannot be empty")
	}

	engine.config = config

	return &engine, nil
}

func checkRegion(region string) error {
	match, _ := regexp.MatchString("(eu|na|fe)", region)
	if !match {
		return errors.New(`Please provide one of: "eu", "na" or "fe"`)
	}
	return nil
}

func (engine *Engine) GetSellerCentralURLForRegion() string {
	switch engine.config.Region {
	case "eu":
		return "url EU to be defined"
	case "fe":
		return "url FE to be defined"
	default:
		// default is na
		return "https://sellercentral.amazon.com"
	}
}

func (engine *Engine) GetLWAURL() (url string) {
	baseURL := engine.GetSellerCentralURLForRegion()
	url = fmt.Sprintf("%s/apps/authorize/consent?application_id=%s", baseURL, engine.config.Credentials.ClientID)

	if engine.config.Beta {
		url = fmt.Sprintf("%s&version=beta", url)
	}

	return url
}

func (engine *Engine) Authenticate(code string) (*OAuthSPAPI, error) {
	resp, err := http.PostForm("https://api.amazon.com/auth/o2/token", url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {engine.config.SPAPICallbackURL},
		"client_id":     {engine.config.Credentials.ClientID},
		"client_secret": {engine.config.Credentials.SPAPIClientSecret}})

	if err != nil {
		return nil, errors.New("Error creating post request")
	}
	defer resp.Body.Close()
	// Se o request deu algum erro alerta
	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		errorStr := buf.String()

		return nil, errors.New(fmt.Sprintf("Fail to generate token. %s", errorStr))
	}

	decoder := json.NewDecoder(resp.Body)
	oauthSPAPI := OAuthSPAPI{}
	err = decoder.Decode(&oauthSPAPI)

	engine.config.Credentials.AccessToken = oauthSPAPI.AccessToken
	engine.config.Credentials.RefreshToken = oauthSPAPI.RefreshToken
	engine.config.Credentials.TokenType = oauthSPAPI.TokenType
	engine.config.Credentials.ExpiresIn = oauthSPAPI.ExpiresIn

	return &oauthSPAPI, nil
}
