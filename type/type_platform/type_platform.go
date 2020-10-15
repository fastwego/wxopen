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

package type_platform

import (
	"encoding/xml"
)

const (
	EventTypeComponentVerifyTicket = "component_verify_ticket"
	EventTypeAuthorized            = "authorized"
	EventTypeUnauthorized          = "unauthorized"
	EventTypeUpdateAuthorized      = "updateauthorized"
)

type Event struct {
	XMLName    xml.Name `xml:"xml"`
	AppId      string
	CreateTime string
	InfoType   string
}

/*
<xml>
<AppId>some_appid</AppId>
<CreateTime>1413192605</CreateTime>
<InfoType>component_verify_ticket</InfoType>
<ComponentVerifyTicket>some_verify_ticket</ComponentVerifyTicket>
</xml>
*/
type EventComponentVerifyTicket struct {
	Event
	ComponentVerifyTicket string
}

/*
授权成功通知
<xml>
  <AppId>第三方平台appid</AppId>
  <CreateTime>1413192760</CreateTime>
  <InfoType>authorized</InfoType>
  <AuthorizerAppid>公众号appid</AuthorizerAppid>
  <AuthorizationCode>授权码</AuthorizationCode>
  <AuthorizationCodeExpiredTime>过期时间</AuthorizationCodeExpiredTime>
  <PreAuthCode>预授权码</PreAuthCode>
<xml>
*/
type EventAuthorized struct {
	Event
	AuthorizerAppid              string
	AuthorizationCode            string
	AuthorizationCodeExpiredTime string
	PreAuthCode                  string
}

/*
取消授权通知
<xml>
  <AppId>第三方平台appid</AppId>
  <CreateTime>1413192760</CreateTime>
  <InfoType>unauthorized</InfoType>
  <AuthorizerAppid>公众号appid</AuthorizerAppid>
</xml>
*/
type EventUnauthorized struct {
	Event
	AuthorizerAppid string
}

/*
授权更新通知
<xml>
  <AppId>第三方平台appid</AppId>
  <CreateTime>1413192760</CreateTime>
  <InfoType>updateauthorized</InfoType>
  <AuthorizerAppid>公众号appid</AuthorizerAppid>
  <AuthorizationCode>授权码</AuthorizationCode>
  <AuthorizationCodeExpiredTime>过期时间</AuthorizationCodeExpiredTime>
  <PreAuthCode>预授权码</PreAuthCode>
<xml>
*/
type EventUpdateAuthorized struct {
	Event
	AuthorizerAppid              string
	AuthorizationCode            string
	AuthorizationCodeExpiredTime string
	PreAuthCode                  string
}
