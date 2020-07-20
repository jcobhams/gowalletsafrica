package gowalletsafrica

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//GetBanks - gets a list of nigerian banks and their code
//Documentation: https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#16f2d271-c546-46a0-b90f-cb6ed06e44b7
func (p *payouts) GetBanks() (Banks, error) {
	banks := Banks{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/transfer/banks/all", p.APIURL), nil)
	if err != nil {
		return banks, err
	}

	resp, err := p.makeRequest(req)
	if err != nil {
		return banks, err
	}
	defer resp.Body.Close()

	var decodedResponseBody []map[string]interface{}
	rawResponseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return banks, err
	}

	err = json.Unmarshal(rawResponseBody, &decodedResponseBody)
	if err != nil {
		return banks, err
	}

	if resp.StatusCode != http.StatusOK {
		return banks, errors.New(fmt.Sprintf("Request Failed - Error Code: %v | Message: %v", resp.StatusCode, resp.Status))
	}

	for _, b := range decodedResponseBody {
		bank := Bank{
			BankCode:     b["BankCode"].(string),
			BankName:     b["BankName"].(string),
			BankSortCode: b["BankSortCode"].(string),
		}
		banks = append(banks, bank)
	}
	return banks, nil
}

//BankDetails - Get transaction details about wallet to bank transfer
//Documentation: https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#bd8ec9a7-2330-4694-a799-0961b0b7bf01
func (p *payouts) BankDetails(transactionReference string) (BankDetail, error) {
	bankDetail := BankDetail{}

	payloadValues := payloadBody{
		"SecretKey":            p.secretKey,
		"TransactionReference": transactionReference,
	}

	payload, err := json.Marshal(payloadValues)
	if err != nil {
		return bankDetail, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/transfer/bank/details", p.APIURL), bytes.NewReader(payload))
	if err != nil {
		return bankDetail, err
	}

	resp, err := p.makeRequest(req)
	if err != nil {
		return bankDetail, err
	}
	defer resp.Body.Close()

	decodedResponseBody, err := p.unmarshallJson(resp.Body)
	if err != nil {
		return bankDetail, err
	}

	if resp.StatusCode != http.StatusOK {
		return bankDetail, errors.New(fmt.Sprintf("Request Failed - Error Code: %v | Message: %v", p.getResponseCode(decodedResponseBody), p.getResponseMessage(decodedResponseBody)))
	}

	bankDetail.Bank = decodedResponseBody["Bank"].(string)
	bankDetail.AccountNumber = decodedResponseBody["AccountNumber"].(string)
	bankDetail.DateTransferred = decodedResponseBody["DateTransferred"].(string)
	bankDetail.Amount = decodedResponseBody["Amount"].(float64)
	bankDetail.RecipientName = decodedResponseBody["RecipientName"].(string)
	bankDetail.ResponseCode = decodedResponseBody["ResponseCode"].(string)

	if decodedResponseBody["SessionId"] != nil {
		bankDetail.SessionId = decodedResponseBody["SessionId"].(string)
	}

	if decodedResponseBody["Message"] != nil {
		bankDetail.Message = decodedResponseBody["Message"].(string)
	}

	return bankDetail, nil
}
