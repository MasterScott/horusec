// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Requests save in this file are exclusive of horusec e2e
package horusec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	accountentities "github.com/ZupIT/horusec/development-kit/pkg/entities/account"
	authEntities "github.com/ZupIT/horusec/development-kit/pkg/entities/auth"
	authDto "github.com/ZupIT/horusec/development-kit/pkg/entities/auth/dto"
	"github.com/stretchr/testify/assert"
)

func CreateAccount(t *testing.T, account *authEntities.Account) {
	fmt.Println("Running test for CreateAccount")
	createAccountResp, err := http.Post("http://127.0.0.1:8006/auth/account/create-account", "text/json", bytes.NewReader(account.ToBytes()))
	assert.NoError(t, err, "create account error mount request")
	assert.Equal(t, http.StatusCreated, createAccountResp.StatusCode, "create account error send request")

	var createAccountResponse map[string]interface{}
	_ = json.NewDecoder(createAccountResp.Body).Decode(&createAccountResponse)
	assert.NoError(t, createAccountResp.Body.Close())
	assert.NotEmpty(t, createAccountResponse["content"])
}

func Login(t *testing.T, credentials *authDto.Credentials) map[string]string {
	fmt.Println("Running test for Login")
	loginResp, err := http.Post(
		"http://127.0.0.1:8006/auth/auth/authenticate",
		"text/json",
		bytes.NewReader(credentials.ToBytes()),
	)
	assert.NoError(t, err, "login, error mount request")
	assert.Equal(t, http.StatusOK, loginResp.StatusCode, "login error send request")

	var loginResponse map[string]map[string]string
	_ = json.NewDecoder(loginResp.Body).Decode(&loginResponse)
	assert.NoError(t, loginResp.Body.Close())
	return loginResponse["content"]
}

func Logout(t *testing.T, bearerToken string) {
	fmt.Println("Running test for Logout")
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8006/auth/account/logout", nil)
	req.Header.Add("X-Horusec-Authorization", bearerToken)
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	assert.NoError(t, err, "logout error mount request")
	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "logout error send request")

	var logoutResponse map[string]map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&logoutResponse)
	assert.NoError(t, resp.Body.Close())
}

func CreateCompanyApplicationAdmin(t *testing.T, bearerToken string, company *accountentities.CompanyApplicationAdmin) (CompanyID string) {
	companyBytes, _ := json.Marshal(company)
	fmt.Println("Running test for CreateCompany")
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8003/account/companies", bytes.NewReader(companyBytes))
	req.Header.Add("X-Horusec-Authorization", bearerToken)
	httpClient := http.Client{}
	createCompanyResp, err := httpClient.Do(req)
	assert.NoError(t, err, "create company error send request")
	assert.Equal(t, http.StatusCreated, createCompanyResp.StatusCode, "create company error check response")
	var createdCompany map[string]map[string]string
	_ = json.NewDecoder(createCompanyResp.Body).Decode(&createdCompany)
	assert.NoError(t, createCompanyResp.Body.Close())
	assert.NotEmpty(t, createdCompany["content"]["companyID"])
	return createdCompany["content"]["companyID"]
}

func UpdateCompany(t *testing.T, bearerToken string, companyID string, company *accountentities.Company) {
	fmt.Println("Running test for UpdateCompany")
	req, _ := http.NewRequest(http.MethodPatch, "http://127.0.0.1:8003/account/companies/"+companyID, bytes.NewReader(company.ToBytes()))
	req.Header.Add("X-Horusec-Authorization", bearerToken)
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	assert.NoError(t, err, "update company error send request")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "update company error check response")
	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, resp.Body.Close())
	assert.NotEmpty(t, body["content"])
}

func ReadAllCompanies(t *testing.T, bearerToken string, isCheckBodyEmpty bool) string {
	fmt.Println("Running test for ReadAllCompanies")
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8003/account/companies", nil)
	req.Header.Add("X-Horusec-Authorization", bearerToken)
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	assert.NoError(t, err, "read all companies error send request")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "read all companies error check response")
	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, resp.Body.Close())
	if isCheckBodyEmpty {
		assert.NotEmpty(t, body["content"])
	}
	content, _ := json.Marshal(body["content"])
	return string(content)
}

func DeleteCompany(t *testing.T, bearerToken, companyID string) {
	fmt.Println("Running test for DeleteCompany")
	req, _ := http.NewRequest(http.MethodDelete, "http://127.0.0.1:8003/account/companies/"+companyID, nil)
	req.Header.Add("X-Horusec-Authorization", bearerToken)
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	assert.NoError(t, err, "delete company error send request")
	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "delete company error check response")
	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, resp.Body.Close())
}
