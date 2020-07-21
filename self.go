package gowalletsafrica

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

//CheckBalance retrieves the wallet balance in provided Currency
//Documentation: https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#b9a4e222-3e51-4ff5-b93c-9a36e87be2f7
func (s *self) CheckBalance(currency Currency) (CheckBalanceResult, error) {
	result := CheckBalanceResult{}
	payloadValues := payloadBody{
		"Currency":  currency,
		"SecretKey": s.secretKey,
	}

	payload, err := json.Marshal(payloadValues)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/self/balance", s.APIURL), bytes.NewReader(payload))
	if err != nil {
		return result, err
	}

	resp, err := s.makeRequest(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	decodedResponseBody, err := s.unmarshallJson(resp.Body)
	if err != nil {
		return result, err
	}

	if resp.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("Request Failed - Error Code: %v | Message: %v", s.getResponseCode(decodedResponseBody), s.getResponseMessage(decodedResponseBody)))
	}

	result.WalletBalance = decodedResponseBody["Data"].(map[string]interface{})["WalletBalance"].(float64)
	result.WalletCurrency = decodedResponseBody["Data"].(map[string]interface{})["WalletCurrency"].(string)

	return result, nil
}

//Transactions gets a list of transactions based on provided args or error
//Documentation: https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#32e58247-5227-4d40-aa37-94c01e4886a7
func (s *self) Transactions(currency Currency, transactionType TransactionType, take, skip int, dateFrom, dateTo string) (Transactions, error) {
	transactions := Transactions{}

	if take < 1 {
		return transactions, errors.New("take cannot be less than 1")
	}

	payloadValues := payloadBody{
		"Currency":        currency,
		"TransactionType": transactionType,
		"Take":            take,
		"Skip":            skip,
		"SecretKey":       s.secretKey,
	}

	if dateFrom != "" {
		_, err := time.Parse(DateFormat, dateFrom)
		if err != nil {
			return transactions, err
		}
		payloadValues["DateFrom"] = dateFrom
	}

	if dateTo != "" {
		_, err := time.Parse(DateFormat, dateTo)
		if err != nil {
			return transactions, err
		}
		payloadValues["DateTo"] = dateTo
	}

	payload, err := json.Marshal(payloadValues)
	if err != nil {
		return transactions, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/self/transactions", s.APIURL), bytes.NewReader(payload))
	if err != nil {
		return transactions, err
	}

	resp, err := s.makeRequest(req)
	if err != nil {
		return transactions, err
	}
	defer resp.Body.Close()

	decodedResponseBody, err := s.unmarshallJson(resp.Body)
	if err != nil {
		return transactions, err
	}

	if resp.StatusCode != http.StatusOK {
		return transactions, errors.New(fmt.Sprintf("Request Failed - Error Code: %v | Message: %v", s.getResponseCode(decodedResponseBody), s.getResponseMessage(decodedResponseBody)))
	}

	for _, t := range decodedResponseBody["Data"].(map[string]interface{})["Transactions"].([]interface{}) {
		transaction := Transaction{
			Amount:          t.(map[string]interface{})["Amount"].(float64),
			Currency:        t.(map[string]interface{})["Currency"].(string),
			Category:        t.(map[string]interface{})["Category"].(string),
			Narration:       t.(map[string]interface{})["Narration"].(string),
			DateTransacted:  t.(map[string]interface{})["DateTransacted"].(string),
			PreviousBalance: t.(map[string]interface{})["PreviousBalance"].(float64),
			NewBalance:      t.(map[string]interface{})["NewBalance"].(float64),
			Type:            t.(map[string]interface{})["Type"].(string),
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

//GetWallets - retrieves a list of all wallets created.
//Documentation: https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#e42d61b9-11a3-4dbe-b95e-8641812b0919
func (s *self) GetWallets() (Wallets, error) {
	wallets := Wallets{}

	payloadValues := payloadBody{
		"SecretKey": s.secretKey,
	}

	payload, err := json.Marshal(payloadValues)
	if err != nil {
		return wallets, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/self/users", s.APIURL), bytes.NewReader(payload))
	if err != nil {
		return wallets, err
	}

	resp, err := s.makeRequest(req)
	if err != nil {
		return wallets, err
	}
	defer resp.Body.Close()

	decodedResponseBody, err := s.unmarshallJson(resp.Body)
	if err != nil {
		return wallets, err
	}

	if resp.StatusCode != http.StatusOK {
		return wallets, errors.New(fmt.Sprintf("Request Failed - Error Code: %v | Message: %v", s.getResponseCode(decodedResponseBody), s.getResponseMessage(decodedResponseBody)))
	}

	for _, w := range decodedResponseBody["Data"].([]interface{}) {
		wallet := Wallet{
			Username:    "",
			BVN:         "",
			City:        "",
			Country:     "",
			DateCreated: w.(map[string]interface{})["DateCreated"].(string),
			Email:       w.(map[string]interface{})["Email"].(string),
			FirstName:   w.(map[string]interface{})["FirstName"].(string),
			LastName:    w.(map[string]interface{})["LastName"].(string),
			PhoneNumber: w.(map[string]interface{})["PhoneNumber"].(string),
		}

		if w.(map[string]interface{})["Username"] != nil {
			wallet.Username = w.(map[string]interface{})["Username"].(string)
		}

		if w.(map[string]interface{})["BVN"] != nil {
			wallet.BVN = w.(map[string]interface{})["BVN"].(string)
		}

		if w.(map[string]interface{})["City"] != nil {
			wallet.City = w.(map[string]interface{})["City"].(string)
		}

		if w.(map[string]interface{})["Country"] != nil {
			wallet.Country = w.(map[string]interface{})["Country"].(string)
		}

		if w.(map[string]interface{})["DateOfBirth"] != nil {
			wallet.DateOfBirth = w.(map[string]interface{})["DateOfBirth"].(string)
		}

		wallets = append(wallets, wallet)
	}
	return wallets, nil
}
