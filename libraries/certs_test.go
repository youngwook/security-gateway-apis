/*******************************************************************************
 * Copyright 2018 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * @author: Tingyu Zeng, Dell
 * @version: 1.0.0
 *******************************************************************************/
package libraries

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRequestor struct {
	SecretSvcBaseURL string
}

func (tr *testRequestor) GetProxyBaseURL() string {
	return "test"
}

func (tr *testRequestor) GetSecretSvcBaseURL() string {
	return tr.SecretSvcBaseURL
}

func (tr *testRequestor) GetHttpClient() *http.Client {
	return &http.Client{}
}

type testCertCfg struct {
	CertPath string
}

func (tc *testCertCfg) GetCertPath() string {
	return tc.CertPath
}

func (tc *testCertCfg) GetTokenPath() string {
	return "test"
}

func TestGetSecret(t *testing.T) {
	path := "../../../test/test-resp-init.json"
	cs := Certs{&testRequestor{}, &testCertCfg{}}
	s, err := cs.getSecret(path)
	if err != nil {
		t.Errorf("failed to parse token file")
		t.Errorf(err.Error())
	}
	if s != "test-token" {
		t.Errorf("incorrect token")
		t.Errorf(s)
	}
}

func TestValidate(t *testing.T) {
	cp := &CertPair{"private-cert", "private-key"}
	cs := Certs{&testRequestor{}, &testCertCfg{}}
	err := cs.validate(cp)
	if err != nil {
		t.Errorf("failed to validate cert collection")
	}
}

func TestRetrieve(t *testing.T) {
	certPath := "testCertPath"
	token := "token"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s instead", r.Method)
		}

		if r.URL.EscapedPath() != fmt.Sprintf("/%s", certPath) {
			t.Errorf("expected request to /%s, got %s instead", certPath, r.URL.EscapedPath())
		}

		if r.Header.Get(VaultToken) != token {
			t.Errorf("expected request header for %s is %s, got %s instead", VaultToken, token, r.Header.Get(VaultToken))
		}
	}))
	defer ts.Close()

	cs := Certs{&testRequestor{ts.URL}, &testCertCfg{certPath}}
	_, err := cs.retrieve(token)
	if err != nil {
		t.Errorf("failed to retrieve cert pair")
		t.Errorf(err.Error())
	}
}
