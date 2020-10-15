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

package wxopen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	WXServerUrl                     = "https://api.weixin.qq.com" // 微信 api 服务器地址
	UserAgent                       = "fastwego/wxopen"
	ErrorComponentAccessTokenExpire = errors.New("component_access_token expire")
	ErrorSystemBusy                 = errors.New("system busy")
)

/*
Client 用于向微信接口发送请求
*/
type Client struct {
	Ctx *Platform
}

// HTTPGet GET 请求
func (client *Client) HTTPGet(uri string) (resp []byte, err error) {
	newUrl, err := client.applyComponentAccessToken(uri)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, WXServerUrl+newUrl, nil)
	if err != nil {
		return
	}

	return client.httpDo(req)
}

//HTTPPost POST 请求
func (client *Client) HTTPPost(uri string, payload io.Reader, contentType string) (resp []byte, err error) {
	newUrl, err := client.applyComponentAccessToken(uri)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, WXServerUrl+newUrl, payload)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", contentType)

	return client.httpDo(req)
}

//httpDo 执行 请求
func (client *Client) httpDo(req *http.Request) (resp []byte, err error) {
	req.Header.Add("User-Agent", UserAgent)

	if client.Ctx.Logger != nil {
		client.Ctx.Logger.Printf("%s %s Headers %v", req.Method, req.URL.String(), req.Header)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	resp, err = responseFilter(response)

	// 发现 component_access_token 过期
	if err == ErrorComponentAccessTokenExpire {

		// 主动 通知 component_access_token 过期
		err = client.Ctx.NoticeComponentAccessTokenExpireHandler(client.Ctx)
		if err != nil {
			return
		}

		// 通知到位后 component_access_token 会被刷新，那么可以 retry 了
		var accessToken string
		accessToken, err = client.Ctx.GetComponentAccessTokenHandler(client.Ctx)
		if err != nil {
			return
		}

		// 换新
		q := req.URL.Query()
		q.Set("component_access_token", accessToken)
		req.URL.RawQuery = q.Encode()

		if client.Ctx.Logger != nil {
			client.Ctx.Logger.Printf("%v retry %s %s Headers %v", ErrorComponentAccessTokenExpire, req.Method, req.URL.String(), req.Header)
		}

		response, err = http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer response.Body.Close()

		resp, err = responseFilter(response)
	}

	// -1 系统繁忙，此时请开发者稍候再试
	// 重试一次
	if err == ErrorSystemBusy {

		if client.Ctx.Logger != nil {
			client.Ctx.Logger.Printf("%v : retry %s %s Headers %v", ErrorSystemBusy, req.Method, req.URL.String(), req.Header)
		}

		response, err = http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer response.Body.Close()

		resp, err = responseFilter(response)
	}

	return
}

/*
在请求地址上附加上 component_access_token
*/
func (client *Client) applyComponentAccessToken(oldUrl string) (newUrl string, err error) {
	accessToken, err := client.Ctx.GetComponentAccessTokenHandler(client.Ctx)
	if err != nil {
		return
	}
	if strings.Contains(oldUrl, "?") {
		newUrl = oldUrl + "&component_access_token=" + accessToken
	} else {
		newUrl = oldUrl + "?component_access_token=" + accessToken
	}
	return
}

/*
筛查微信 api 服务器响应，判断以下错误：

- http 状态码 不为 200

- 接口响应错误码 errcode 不为 0
*/
func responseFilter(response *http.Response) (resp []byte, err error) {
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("Status %s", response.Status)
		return
	}

	resp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	errorResponse := struct {
		Errcode int64  `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(resp, &errorResponse)
	if err != nil {
		return
	}

	// 40001(覆盖刷新超过5min后，使用旧 access_token 报错) 获取 access_token 时 AppSecret 错误，或者 access_token 无效。请开发者认真比对 AppSecret 的正确性，或查看是否正在为恰当的公众号调用接口
	// 42001(超过 7200s 后 报错) - access_token 超时，请检查 access_token 的有效期，请参考基础支持 - 获取 access_token 中，对 access_token 的详细机制说明
	if errorResponse.Errcode == 42001 || errorResponse.Errcode == 40001 {
		err = ErrorComponentAccessTokenExpire
		return
	}

	//  -1	系统繁忙，此时请开发者稍候再试
	if errorResponse.Errcode == -1 {
		err = ErrorSystemBusy
		return
	}

	if errorResponse.Errcode != 0 {
		err = errors.New(string(resp))
		return
	}
	return
}
