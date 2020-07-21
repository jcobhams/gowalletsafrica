package gowalletsafrica

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (w *wallets) Generate(currency Currency, firstName, lastName, email, dateOfBirth string) (Wallet, error) {
	wallet := Wallet{}

	payloadValues := payloadBody{
		"SecretKey":   w.secretKey,
		"Currency":    currency,
		"FirstName":   firstName,
		"LastName":    lastName,
		"Email":       email,
		"DateOfBirth": dateOfBirth,
	}

	payload, err := json.Marshal(payloadValues)
	if err != nil {
		return wallet, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/wallet/generate", w.APIURL), bytes.NewReader(payload))
	if err != nil {
		return wallet, err
	}

	resp, err := w.makeRequest(req)
	if err != nil {
		return wallet, err
	}
	defer resp.Body.Close()

	decodedResponseBody, err := w.unmarshallJson(resp.Body)
	if err != nil {
		return wallet, err
	}

	if resp.StatusCode != http.StatusOK {
		return wallet, errors.New(fmt.Sprintf("Request Failed - Error Code: %v | Message: %v", w.getResponseCode(decodedResponseBody), w.getResponseMessage(decodedResponseBody)))
	}

	data := decodedResponseBody["Data"].(map[string]interface{})
	wallet.FirstName = data["FirstName"].(string)
	wallet.LastName = data["LastName"].(string)
	wallet.Email = data["Email"].(string)
	wallet.PhoneNumber = data["PhoneNumber"].(string)
	if data["BVN"] != nil {
		wallet.BVN = data["BVN"].(string)
	}
	wallet.Password = data["Password"].(string)
	wallet.DateOfBirth = data["DateOfBirth"].(string)
	wallet.DateSignedup = data["DateSignedup"].(string)
	wallet.AccountNo = data["AccountNo"].(string)
	wallet.Bank = data["Bank"].(string)
	wallet.AccountName = data["AccountName"].(string)
	wallet.AvailableBalance = data["AvailableBalance"].(float64)

	return wallet, nil
}
