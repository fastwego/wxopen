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

// Package oauth 代公众号发起网页授权
package oauth

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/fastwego/wxopen"
)

const (
	apiGetAccessToken     = "/sns/oauth2/component/access_token"
	apiRefreshAccessToken = "/sns/oauth2/component/refresh_token"
	apiGetUserInfo        = "/sns/userinfo"
)

/*
获取用户授权跳转链接

在确保微信公众账号拥有授权作用域（scope 参数）的权限的前提下（一般而言，已微信认证的服务号拥有 snsapi_base 和 snsapi_userinfo），使用微信客户端打开以下链接（严格按照以下格式，包括顺序和大小写，并请将参数替换为实际内容）

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Official_Accounts/official_account_website_authorization.html

GET https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE&component_appid=component_appid#wechat_redirect

*/
func GetAuthorizeUrl(ctx *wxopen.Platform, appid string, redirect_uri string, scope string, state string) (redirectUri string, err error) {
	uriTpl := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s&component_appid=%s#wechat_redirect"
	redirectUri = fmt.Sprintf(uriTpl, appid, url.QueryEscape(redirect_uri), scope, state, ctx.Config.AppId)
	return
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

/*
通过code换取网页授权access_token

获取第一步的 code 后，请求以下链接获取 access_token 需要注意的是，由于安全方面的考虑，对访问该链接的客户端有 IP 白名单的要求

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Official_Accounts/official_account_website_authorization.html

GET https://api.weixin.qq.com/sns/oauth2/component/access_token?appid=APPID&code=CODE&grant_type=authorization_code&component_appid=COMPONENT_APPID&component_access_token=COMPONENT_ACCESS_TOKEN
*/
func GetAccessToken(ctx *wxopen.Platform, appid string, code string) (accessToken AccessToken, err error) {

	params := url.Values{}
	params.Add("appid", appid)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")
	params.Add("component_appid", ctx.Config.AppId)

	resp, err := ctx.Client.HTTPGet(apiGetAccessToken + "?" + params.Encode())
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &accessToken)
	if err != nil {
		return
	}
	return
}

/*
刷新access_token

由于 access_token 拥有较短的有效期，当 access_token 超时后，可以使用 refresh_token 进行刷新，refresh_token 拥有较长的有效期（30 天），当 refresh_token 失效的后，需要用户重新授权

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Official_Accounts/official_account_website_authorization.html

GET https://api.weixin.qq.com/sns/oauth2/component/refresh_token?appid=APPID&grant_type=refresh_token&component_appid=COMPONENT_APPID&component_access_token=COMPONENT_ACCESS_TOKEN&refresh_token=REFRESH_TOKEN
*/
func RefreshAccessToken(ctx *wxopen.Platform, appid string, refresh_token string) (accessToken AccessToken, err error) {
	params := url.Values{}
	params.Add("appid", appid)
	params.Add("refresh_token", refresh_token)
	params.Add("grant_type", "refresh_token")
	params.Add("component_appid", ctx.Config.AppId)

	resp, err := ctx.Client.HTTPGet(apiRefreshAccessToken + "?" + params.Encode())
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &accessToken)
	if err != nil {
		return
	}
	return

}

type UserInfo struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int64    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

/*
拉取用户信息

如果网页授权作用域为snsapi_userinfo，则此时开发者可以通过access_token和openid拉取用户信息了

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Official_Accounts/official_account_website_authorization.html

GET https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
*/
func GetUserInfo(ctx *wxopen.Platform, access_token string, openid string) (userInfo UserInfo, err error) {

	params := url.Values{}
	params.Add("access_token", access_token)
	params.Add("openid", openid)
	params.Add("lang", "zh_CN")

	resp, err := ctx.Client.HTTPGet(apiGetUserInfo + "?" + params.Encode())
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &userInfo)
	if err != nil {
		return
	}
	return
}
