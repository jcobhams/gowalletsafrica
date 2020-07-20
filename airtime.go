package gowalletsafrica

import (
	"errors"
	"fmt"
	"net/http"
)

//GetProviders - returns a list of Network Providers
//Documentation: https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#c397d163-a6ff-4bfa-92ce-b75a2a3e9cbd
func (a *airtime) GetProviders() (AirtimeProviders, error) {
	providers := AirtimeProviders{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/bills/airtime/providers", a.APIURL), nil)
	if err != nil {
		return providers, err
	}

	resp, err := a.makeRequest(req)
	if err != nil {
		return providers, err
	}
	defer resp.Body.Close()

	decodedResponseBody, err := a.unmarshallJson(resp.Body)
	if err != nil {
		return providers, err
	}

	if resp.StatusCode != http.StatusOK {
		return providers, errors.New(fmt.Sprintf("Request Failed - Error Code: %v | Message: %v", a.getResponseCode(decodedResponseBody), a.getResponseMessage(decodedResponseBody)))
	}

	for _, p := range decodedResponseBody["Providers"].([]interface{}) {
		provider := AirtimeProvider{
			Code: p.(map[string]interface{})["Code"].(string),
			Name: p.(map[string]interface{})["Name"].(string),
		}
		providers = append(providers, provider)
	}

	return providers, nil
}
