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
	"net/http"
)

type Requestor interface {
	GetProxyBaseURL() string
	GetSecretSvcBaseURL() string
	GetHttpClient() *http.Client
}

type EdgeXRequestor struct {
	ProxyBaseURL     string
	SecretSvcBaseURL string
	Client           *http.Client
}

func (eq *EdgeXRequestor) GetProxyBaseURL() string {
	return eq.ProxyBaseURL
}

func (eq *EdgeXRequestor) GetSecretSvcBaseURL() string {
	return eq.SecretSvcBaseURL
}

func (eq *EdgeXRequestor) GetHttpClient() *http.Client {
	return eq.Client
}
