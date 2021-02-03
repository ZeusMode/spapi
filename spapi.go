package spapi

import (
	"errors"
	"fmt"
	"regexp"
)

type Engine struct {
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
}

func New(e Engine) (engine *Engine, err error) {

	err = checkRegion(e.Region)
	if err != nil {
		return nil, err
	}

	if e.Credentials.ClientID == "" {
		return nil, errors.New("Credentials.ClientID cannot be empty")
	}

	if e.Credentials.SPAPIClientID == "" {
		return nil, errors.New("Credentials.SPAPIClientID cannot be empty")
	}

	if e.Credentials.SPAPIClientSecret == "" {
		return nil, errors.New("Credentials.SPAPIClientSecret cannot be empty")
	}

	if e.Credentials.SPAPICallbackURL == "" {
		return nil, errors.New("Credentials.SPAPICallbackURL cannot be empty")
	}

	if e.Credentials.AWSAccessKeyID == "" {
		return nil, errors.New("Credentials.AWSAccessKeyID cannot be empty")
	}

	if e.Credentials.AWSAccessKey == "" {
		return nil, errors.New("Credentials.AWSAccessKey cannot be empty")
	}

	return &e, nil
}

func checkRegion(region string) error {
	match, _ := regexp.MatchString("(eu|na|fe)", region)
	if !match {
		return errors.New(`Please provide one of: "eu", "na" or "fe"`)
	}
	return nil
}

func GetSellerCentralURLForRegion(engine *Engine) string {
	switch engine.Region {
	case "eu":
		return "url EU to be defined"
	case "fe":
		return "url FE to be defined"
	default:
		// default is na
		return "https://sellercentral.amazon.com"
	}
}

func GetLWAURL(engine *Engine) (url string) {
	baseURL := GetSellerCentralURLForRegion(engine)
	url = fmt.Sprintf("%s/apps/authorize/consent?application_id=%s", baseURL, engine.ClientID)

	if engine.Beta {
		url = fmt.Sprintf("%s&version=beta", url)
	}

	return url
}
