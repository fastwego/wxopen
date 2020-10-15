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

// Package auth 开放平台
package auth

import (
	"bytes"
	"net/url"

	"github.com/fastwego/wxopen"
)

const (
	apiCreatePreauthCode            = "/cgi-bin/component/api_create_preauthcode"
	apiGetAuthorizationRedirectUri  = "/cgi-bin/componentloginpage"
	apiGetAuthorizationRedirectUri2 = "/safe/bindcomponent"
	apiApiQueryAuth                 = "/cgi-bin/component/api_query_auth"
	apiApiAuthorizerToken           = "/cgi-bin/component/api_authorizer_token"
	apiApiGetAuthorizerInfo         = "/cgi-bin/component/api_get_authorizer_info"
	apiApiGetAuthorizerOption       = "/cgi-bin/component/api_get_authorizer_option"
	apiApiSetAuthorizerOption       = "/cgi-bin/component/api_set_authorizer_option"
	apiApiGetAuthorizerList         = "/cgi-bin/component/api_get_authorizer_list"
)

/*
获取 预授权码

预授权码（pre_auth_code）是第三方平台方实现授权托管的必备信息，每个预授权码有效期为 10 分钟。需要先获取令牌才能调用

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/pre_auth_code.html

POST https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=COMPONENT_ACCESS_TOKEN
*/
func CreatePreauthCode(ctx *wxopen.Platform, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiCreatePreauthCode, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
方式一：授权注册页面扫码授权

第三方平台方可以在自己的网站中放置“微信公众号授权”或者“小程序授权”的入口，或生成授权链接放置在移动网页中，引导公众号和小程序管理员进入授权页。

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Authorization_Process_Technical_Description.html

GET https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=xxxx&pre_auth_code=xxxxx&redirect_uri=xxxx&auth_type=xxx
*/
func GetAuthorizationRedirectUri(params url.Values) (uri string) {
	return "https://mp.weixin.qq.com/cgi-bin/componentloginpage?" + params.Encode()
}

/*
方式二：点击移动端链接快速授权

第三方平台方可以生成授权链接，将链接通过移动端直接发给授权管理员，管理员确认后即授权成功

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Authorization_Process_Technical_Description.html

GET https://mp.weixin.qq.com/safe/bindcomponent?action=bindcomponent&auth_type=3&no_scan=1&component_appid=xxxx&pre_auth_code=xxxxx&redirect_uri=xxxx&auth_type=xxx&biz_appid=xxxx#wechat_redirect
*/
func GetAuthorizationRedirectUri2(params url.Values) (uri string) {
	return "https://mp.weixin.qq.com/safe/bindcomponent?" + params.Encode() + "#wechat_redirect"
}

/*
使用授权码获取授权信息

由当用户在第三方平台授权页中完成授权流程后，第三方平台开发者可以在回调 URI 中通过 URL 参数获取授权码。使用以下接口可以换取公众号/小程序的授权信息。建议保存授权信息中的刷新令牌（authorizer_refresh_token）

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/authorization_info.html

POST https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=COMPONENT_ACCESS_TOKEN
*/
func ApiQueryAuth(ctx *wxopen.Platform, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiApiQueryAuth, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
获取/刷新接口调用令牌

在公众号/小程序接口调用令牌（authorizer_access_token）失效时，可以使用刷新令牌（authorizer_refresh_token）获取新的接口调用令牌

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/api_authorizer_token.html

POST https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=COMPONENT_ACCESS_TOKEN
*/
func ApiAuthorizerToken(ctx *wxopen.Platform, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiApiAuthorizerToken, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
获取授权方的帐号基本信息

该 API 用于获取授权方的基本信息，包括头像、昵称、帐号类型、认证类型、微信号、原始ID和二维码图片URL

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/api_get_authorizer_info.html

POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?component_access_token=COMPONENT_ACCESS_TOKEN
*/
func ApiGetAuthorizerInfo(ctx *wxopen.Platform, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiApiGetAuthorizerInfo, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
获取授权方选项信息

本 API 用于获取授权方的公众号/小程序的选项设置信息，如：地理位置上报，语音识别开关，多客服开关

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/api_get_authorizer_option.html

POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_option?component_access_token=COMPONENT_ACCESS_TOKEN
*/
func ApiGetAuthorizerOption(ctx *wxopen.Platform, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiApiGetAuthorizerOption, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
设置授权方选项信息

本 API 用于设置授权方的公众号/小程序的选项信息，如：地理位置上报，语音识别开关，多客服开关

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/api_set_authorizer_option.html

POST https://api.weixin.qq.com/cgi-bin/component/api_set_authorizer_option?component_access_token=COMPONENT_ACCESS_TOKEN
*/
func ApiSetAuthorizerOption(ctx *wxopen.Platform, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiApiSetAuthorizerOption, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
拉取所有已授权的帐号信息

使用本 API 拉取当前所有已授权的帐号基本信息

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/api_get_authorizer_list.html

POST https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_list?component_access_token=COMPONENT_ACCESS_TOKEN
*/
func ApiGetAuthorizerList(ctx *wxopen.Platform, payload []byte) (resp []byte, err error) {
	return ctx.Client.HTTPPost(apiApiGetAuthorizerList, bytes.NewReader(payload), "application/json;charset=utf-8")
}
