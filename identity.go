package gowalletsafrica

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//ResolveBVN - Get information about the provided BVN
//Documentation: https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#86ebd12e-c0e7-4529-86ea-9ed5f6993272
func (i *identity) ResolveBVN(bvn string) (ResolveBVN, error) {
	result := ResolveBVN{}
	if bvn == "" {
		return result, errors.New("BVN number is required")
	}

	payloadValues := payloadBody{
		"BVN":       bvn,
		"SecretKey": i.secretKey,
	}

	payload, err := json.Marshal(payloadValues)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%v/account/resolvebvn", i.APIURL), bytes.NewReader(payload))
	if err != nil {
		return result, err
	}

	resp, err := i.makeRequest(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	decodedResponseBody, err := i.unmarshallJson(resp.Body)
	if err != nil {
		return result, err
	}

	if resp.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("Request Failed - Error Code: %v | Message: %v", i.getResponseCode(decodedResponseBody), i.getResponseMessage(decodedResponseBody)))
	}

	result.FirstName = decodedResponseBody["FirstName"].(string)
	result.LastName = decodedResponseBody["LastName"].(string)
	result.Email = decodedResponseBody["Email"].(string)
	result.PhoneNumber = decodedResponseBody["PhoneNumber"].(string)
	result.BVN = decodedResponseBody["BVN"].(string)
	result.DateOfBirth = decodedResponseBody["DateOfBirth"].(string)
	if decodedResponseBody["MiddleName"] != nil {
		result.MiddleName = decodedResponseBody["MiddleName"].(string)
	}
	if decodedResponseBody["EnrollmentBank"] != nil {
		result.EnrollmentBank = decodedResponseBody["EnrollmentBank"].(string)
	}
	if decodedResponseBody["EnrollmentBranch"] != nil {
		result.EnrollmentBranch = decodedResponseBody["EnrollmentBranch"].(string)
	}
	if decodedResponseBody["Gender"] != nil {
		result.Gender = decodedResponseBody["Gender"].(string)
	}
	if decodedResponseBody["LevelOfAccount"] != nil {
		result.LevelOfAccount = decodedResponseBody["LevelOfAccount"].(string)
	}
	if decodedResponseBody["LgaOfOrigin"] != nil {
		result.LgaOfOrigin = decodedResponseBody["LgaOfOrigin"].(string)
	}
	if decodedResponseBody["LgaOfResidence"] != nil {
		result.LgaOfResidence = decodedResponseBody["LgaOfResidence"].(string)
	}
	if decodedResponseBody["MaritalStatus"] != nil {
		result.MaritalStatus = decodedResponseBody["MaritalStatus"].(string)
	}
	if decodedResponseBody["NameOnCard"] != nil {
		result.NameOnCard = decodedResponseBody["NameOnCard"].(string)
	}
	if decodedResponseBody["Nationality"] != nil {
		result.Nationality = decodedResponseBody["Nationality"].(string)
	}
	if decodedResponseBody["StateOfOrigin"] != nil {
		result.StateOfOrigin = decodedResponseBody["StateOfOrigin"].(string)
	}
	if decodedResponseBody["StateOfResidence"] != nil {
		result.StateOfResidence = decodedResponseBody["StateOfResidence"].(string)
	}
	if decodedResponseBody["Title"] != nil {
		result.Title = decodedResponseBody["Title"].(string)
	}
	if decodedResponseBody["WatchListed"] != nil {
		result.WatchListed = decodedResponseBody["WatchListed"].(string)
	}
	if decodedResponseBody["Picture"] != nil {
		result.Picture = decodedResponseBody["Picture"].(string)
	}

	return result, nil
}

//ResolveBVNDetails - Get full BVN details. It's an alias of the first as they share the same result set.
//Documentation: https://documenter.getpostman.com/view/10058163/SWLk4RPL?version=latest#8fabbdbd-2e82-4138-9d9b-566bad20d55b
func (i *identity) ResolveBVNDetails(bvn string) (ResolveBVN, error) {
	return i.ResolveBVN(bvn)
}
