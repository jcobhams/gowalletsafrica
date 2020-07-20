package gowalletsafrica

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	EnvSandbox string = "sandbox"
	EnvLive    string = "live"

	SandBoxPublicKey string = "uvjqzm5xl6bw"
	SandBoxSecretKey string = "hfucj5jatq8h"

	APIBaseUrlSandbox string = "https://sandbox.wallets.africa"
	APIBaseUrlLive    string = "https://api.wallets.africa"

	RequestTimeout time.Duration = 5 * time.Second

	CurrencyNigeria Currency = "NGN"
	CurrencyUSA     Currency = "USD"
	CurrencyGhana   Currency = "GHS"
	CurrencyKenya   Currency = "KES"

	TransactionTypeCredit TransactionType = 1
	TransactionTypeDebit  TransactionType = 2
	TransactionTypeAll    TransactionType = 3
)

var DefaultConfig = Config{
	Environment:    EnvSandbox,
	PublicKey:      SandBoxPublicKey,
	SecretKey:      SandBoxSecretKey,
	RequestTimeout: RequestTimeout,
}

//New create a new instance of the WalletsAfrica struct based on provided config.
//Returns a pointer to the struct and nil error if successful or a nil pointer and an error
func New(config Config) (*WalletsAfrica, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	base := newBase(config)

	wa := &WalletsAfrica{
		Self:     &self{base},
		Wallets:  &wallets{base},
		Payouts:  &payouts{base},
		Airtime:  &airtime{base},
		Identity: &identity{base},
	}
	return wa, nil
}

//validateConfig checks the provided config to ensure it's well formed
func validateConfig(config Config) error {
	if config.Environment != EnvSandbox && config.Environment != EnvLive {
		return errors.New(fmt.Sprintf("malformed config - provided enviroment is not supported. - Only %v or %v is allowed", EnvLive, EnvSandbox))
	}

	if config.PublicKey == "" {
		return errors.New("malformed config - public key if required")
	}

	if config.SecretKey == "" {
		return errors.New("malformed config - secret key if required")
	}

	if config.Environment == EnvLive && config.PublicKey == SandBoxPublicKey {
		return errors.New("malformed config - using sandbox public key in live mode not permitted")
	}

	if config.Environment == EnvLive && config.SecretKey == SandBoxSecretKey {
		return errors.New("malformed config - using sandbox secret key in live mode not permitted")
	}
	return nil
}

func newBase(config Config) *base {
	b := &base{
		HTTPClient: &http.Client{Timeout: config.RequestTimeout},
		secretKey:  config.SecretKey,
		publicKey:  config.PublicKey,
	}

	switch config.Environment {
	case EnvSandbox:
		b.APIURL = APIBaseUrlSandbox
	case EnvLive:
		b.APIURL = APIBaseUrlLive
	}

	return b
}

func (b *base) makeRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", b.publicKey))

	resp, err := b.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *base) unmarshallJson(body io.Reader) (responseBody, error) {
	responseBody := make(responseBody)
	rawResponseBody, err := ioutil.ReadAll(body)
	if err != nil {
		return responseBody, err
	}

	err = json.Unmarshal(rawResponseBody, &responseBody)
	if err != nil {
		return responseBody, err
	}
	return responseBody, nil
}

func (b *base) getResponseCode(body responseBody) string {
	if code, ok := body["Response"].(map[string]interface{})["ResponseCode"]; ok {
		return code.(string)
	}
	return ""
}

func (b *base) getResponseMessage(body responseBody) string {
	if msg, ok := body["Response"].(map[string]interface{})["Message"]; ok {
		return msg.(string)
	}
	return ""
}
