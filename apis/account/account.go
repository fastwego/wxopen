// Copyright 2020 FastWeGo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package account 开放平台-账号管理
package account

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/fastwego/wxopen"
)

const (
	apiCreate = "/cgi-bin/open/create"
	apiBind   = "/cgi-bin/open/bind"
	apiUnbind = "/cgi-bin/open/unbind"
	apiGet    = "/cgi-bin/open/get"
)

/*
创建开放平台帐号并绑定公众号/小程序

该 API 用于创建一个开放平台帐号，并将一个尚未绑定开放平台帐号的公众号/小程序绑定至该开放平台帐号上。新创建的开放平台帐号的主体信息将设置为与之绑定的公众号或小程序的主体

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/account/create.html

POST https://api.weixin.qq.com/cgi-bin/open/create?access_token=ACCESS_TOKEN
*/
func Create(authorizer_access_token string, payload []byte) (resp []byte, err error) {
	return httpPost(apiCreate+"?access_token="+authorizer_access_token, payload)
}

/*
将公众号/小程序绑定到开放平台帐号下

该 API 用于将一个尚未绑定开放平台帐号的公众号或小程序绑定至指定开放平台帐号上。二者须主体相同。

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/account/bind.html

POST https://api.weixin.qq.com/cgi-bin/open/bind?access_token=xxxx
*/
func Bind(authorizer_access_token string, payload []byte) (resp []byte, err error) {
	return httpPost(apiBind+"?access_token="+authorizer_access_token, payload)
}

/*
将公众号/小程序从开放平台帐号下解绑

该 API 用于将一个公众号或小程序与指定开放平台帐号解绑。开发者须确认所指定帐号与当前该公众号或小程序所绑定的开放平台帐号一致

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/account/unbind.html

POST https://api.weixin.qq.com/cgi-bin/open/unbind?access_token=ACCESS_TOKEN
*/
func Unbind(authorizer_access_token string, payload []byte) (resp []byte, err error) {
	return httpPost(apiUnbind+"?access_token="+authorizer_access_token, payload)
}

/*
获取公众号/小程序所绑定的开放平台帐号

该 API 用于获取公众号或小程序所绑定的开放平台帐号

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/account/get.html

POST https://api.weixin.qq.com/cgi-bin/open/get?access_token=ACCESS_TOKEN
*/
func Get(authorizer_access_token string, payload []byte) (resp []byte, err error) {
	return httpPost(apiGet+"?access_token="+authorizer_access_token, payload)
}

func httpPost(api string, payload []byte) (resp []byte, err error) {
	response, err := http.Post(wxopen.WXServerUrl+api, "application/json;charset=utf-8", bytes.NewReader(payload))
	if err != nil {
		return
	}
	defer response.Body.Close()
	resp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}
