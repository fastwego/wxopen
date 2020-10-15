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

/*
微信开放平台 SDK

See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Third_party_platform_appid.html
*/
package wxopen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fastwego/miniprogram"

	"github.com/fastwego/offiaccount"

	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/file"
)

// GetComponentAccessTokenFunc 获取 component_access_token 方法接口
type GetComponentAccessTokenFunc func(platform *Platform) (accessToken string, err error)

// NoticeComponentAccessTokenExpireFunc 通知中控 刷新 component_access_token
type NoticeComponentAccessTokenExpireFunc func(platform *Platform) (err error)

// GetComponentVerifyTicketFunc 获取 component_verify_ticket 方法接口
type GetComponentVerifyTicketFunc func(platform *Platform) (ticket string, err error)

// ReceiveComponentVerifyTicketFunc 接收 component_verify_ticket 方法接口
type ReceiveComponentVerifyTicketFunc func(platform *Platform, ticket string) (err error)

// GetAuthorizerAccessTokenFunc 获取 AuthorizerAccessToken 方法接口
type GetAuthorizerAccessTokenFunc func(platform *Platform, appid string) (authorizerAccessToken string, err error)

// NoticeAuthorizerAccessTokenExpireFunc 通知刷新 AuthorizerAccessToken 方法接口
type NoticeAuthorizerAccessTokenExpireFunc func(platform *Platform, appid string) (err error)

/*
PlatformConfig 平台 配置
*/
type PlatformConfig struct {
	AppId     string
	AppSecret string
	Token     string
	AesKey    string
}

/*
Platform 平台实例
*/
type Platform struct {
	Config PlatformConfig
	Cache  cachego.Cache
	Client Client
	Server Server
	Logger *log.Logger

	GetComponentAccessTokenHandler          GetComponentAccessTokenFunc
	NoticeComponentAccessTokenExpireHandler NoticeComponentAccessTokenExpireFunc

	GetComponentVerifyTicketHandler     GetComponentVerifyTicketFunc
	ReceiveComponentVerifyTicketHandler ReceiveComponentVerifyTicketFunc

	GetAuthorizerAccessTokenHandler          GetAuthorizerAccessTokenFunc
	NoticeAuthorizerAccessTokenExpireHandler NoticeAuthorizerAccessTokenExpireFunc
}

/*
创建 平台 实例
*/
func NewPlatform(config PlatformConfig) (platform *Platform) {
	instance := Platform{
		Config: config,
		Cache:  file.New(os.TempDir()),

		GetComponentAccessTokenHandler:          GetComponentAccessToken,
		NoticeComponentAccessTokenExpireHandler: NoticeComponentAccessTokenExpire,

		GetComponentVerifyTicketHandler:     GetComponentVerifyTicket,
		ReceiveComponentVerifyTicketHandler: ReceiveComponentVerifyTicket,

		GetAuthorizerAccessTokenHandler:          GetAuthorizerAccessToken,
		NoticeAuthorizerAccessTokenExpireHandler: NoticeAuthorizerAccessTokenExpire,
	}

	instance.Client = Client{Ctx: &instance}
	instance.Server = Server{Ctx: &instance}

	instance.Logger = log.New(os.Stdout, "[fastwego/wxopen] ", log.LstdFlags|log.Llongfile)

	return &instance
}

/*
创建公众号实例
*/
func (platform *Platform) NewOffiAccount(appid string) (offiAccount *offiaccount.OffiAccount, err error) {
	offiAccount = offiaccount.New(offiaccount.Config{
		Appid: appid,
	})

	offiAccount.AccessToken.GetAccessTokenHandler = func(ctx *offiaccount.OffiAccount) (accessToken string, err error) {
		return platform.GetAuthorizerAccessTokenHandler(platform, ctx.Config.Appid)
	}

	offiAccount.AccessToken.NoticeAccessTokenExpireHandler = func(ctx *offiaccount.OffiAccount) (err error) {
		return platform.NoticeAuthorizerAccessTokenExpireHandler(platform, ctx.Config.Appid)
	}

	return
}

/*
创建 小程序 实例
*/
func (platform *Platform) NewMiniprogram(appid string) (mini *miniprogram.Miniprogram, err error) {
	mini = miniprogram.New(miniprogram.Config{
		Appid: appid,
	})

	mini.AccessToken.GetAccessTokenHandler = func(ctx *miniprogram.Miniprogram) (accessToken string, err error) {
		return platform.GetAuthorizerAccessTokenHandler(platform, ctx.Config.Appid)
	}

	mini.AccessToken.NoticeAccessTokenExpireHandler = func(ctx *miniprogram.Miniprogram) (err error) {
		return platform.NoticeAuthorizerAccessTokenExpireHandler(platform, ctx.Config.Appid)
	}

	return
}

/*
GetAuthorizerAccessToken 获取 authorizer_access_token

框架默认将 authorizer_access_token 缓存在本地

实际业务 建议 存储到数据库
*/
func GetAuthorizerAccessToken(platform *Platform, appid string) (accessToken string, err error) {
	accessToken, err = platform.Cache.Fetch("authorizer_access_token:" + appid)
	return
}

// 防止多个 goroutine 并发刷新冲突
var noticeAuthorizerAccessTokenExpireLock sync.Mutex

/*
NoticeAuthorizerAccessTokenExpire 通知 authorizer_access_token 过期

框架默认将 authorizer_access_token 缓存在本地

实际业务 建议 存储到数据库
*/
func NoticeAuthorizerAccessTokenExpire(platform *Platform, appid string) (err error) {

	noticeAuthorizerAccessTokenExpireLock.Lock()
	defer noticeAuthorizerAccessTokenExpireLock.Unlock()

	authorizer_refresh_token, err := platform.Cache.Fetch("authorizer_refresh_token:" + appid)
	if err != nil {
		return
	}

	// 刷新
	params := struct {
		ComponentAppid         string `json:"component_appid"`
		AuthorizerAppid        string `json:"authorizer_appid"`
		AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
	}{
		ComponentAppid:         platform.Config.AppId,
		AuthorizerAppid:        appid,
		AuthorizerRefreshToken: authorizer_refresh_token,
	}

	payload, err := json.Marshal(params)
	apiAuthorizerToken, err := platform.Client.HTTPPost("/cgi-bin/component/api_authorizer_token", bytes.NewReader(payload), "application/json;charset=utf-8")
	if err != nil {
		return
	}
	apiAuthorizerTokenResp := struct {
		AuthorizerAccessToken  string `json:"authorizer_access_token"`
		ExpiresIn              int    `json:"expires_in"`
		AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
	}{}

	err = json.Unmarshal(apiAuthorizerToken, &apiAuthorizerTokenResp)
	if err != nil {
		return
	}

	err = platform.Cache.Save("authorizer_access_token:"+appid, apiAuthorizerTokenResp.AuthorizerAccessToken, 0)
	if err != nil {
		return
	}

	err = platform.Cache.Save("authorizer_refresh_token:"+appid, apiAuthorizerTokenResp.AuthorizerRefreshToken, 0)
	if err != nil {
		return
	}

	return
}

// 防止多个 goroutine 并发刷新冲突
var refreshComponentAccessTokenLock sync.Mutex

/*
从 公众号实例 的 ComponentAccessToken 管理器 获取 access_token

如果没有 access_token 或者 已过期，那么刷新

获得新的 access_token 后 过期时间设置为 0.9 * expiresIn 提供一定冗余
*/
func GetComponentAccessToken(ctx *Platform) (accessToken string, err error) {
	accessToken, err = ctx.Cache.Fetch(ctx.Config.AppId)
	if accessToken != "" {
		return
	}

	refreshComponentAccessTokenLock.Lock()
	defer refreshComponentAccessTokenLock.Unlock()

	accessToken, err = ctx.Cache.Fetch(ctx.Config.AppId)
	if accessToken != "" {
		return
	}

	ticket, err := ctx.GetComponentVerifyTicketHandler(ctx)
	if err != nil {
		return
	}
	accessToken, expiresIn, err := refreshComponentAccessToken(ctx.Config.AppId, ctx.Config.AppSecret, ticket)
	if err != nil {
		return
	}

	// 本地缓存 access_token
	d := time.Duration(expiresIn) * time.Second
	_ = ctx.Cache.Save(ctx.Config.AppId, accessToken, d)

	if ctx.Logger != nil {
		ctx.Logger.Printf("%s %s %d\n", "refreshComponentAccessToken", accessToken, expiresIn)
	}

	return
}

/*
NoticeComponentAccessTokenExpire 只需将本地存储的 access_token 删除，即完成了 access_token 已过期的 主动通知

retry 请求的时候，会发现本地没有 access_token ，从而触发refresh
*/
func NoticeComponentAccessTokenExpire(ctx *Platform) (err error) {
	if ctx.Logger != nil {
		ctx.Logger.Println("NoticeComponentAccessTokenExpire")
	}

	err = ctx.Cache.Delete(ctx.Config.AppId)
	return
}

/*
从微信服务器获取新的 ComponentAccessToken

See: https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
*/
func refreshComponentAccessToken(appid string, secret string, ticket string) (accessToken string, expiresIn int, err error) {
	type Params struct {
		ComponentAppid        string `json:"component_appid"`
		ComponentAppsecret    string `json:"component_appsecret"`
		ComponentVerifyTicket string `json:"component_verify_ticket"`
	}
	params := Params{
		ComponentAppid:        appid,
		ComponentAppsecret:    secret,
		ComponentVerifyTicket: ticket,
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return
	}

	/**
	POST 数据示例：
	{
	  "component_appid":  "appid_value" ,
	  "component_appsecret":  "appsecret_value",
	  "component_verify_ticket": "ticket_value"
	}
	*/
	url := WXServerUrl + "/cgi-bin/component/api_component_token"

	response, err := http.Post(url, "application/json;charset=utf-8", bytes.NewReader(payload))
	if err != nil {
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET %s RETURN %s", url, response.Status)
		return
	}

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	/**
	返回结果示例：
	{
	  "component_access_token": "61W3mEpU66027wgNZ_MhGHNQDHnFATkDa9-2llqrMBjUwxRSNPbVsMmyD-yq8wZETSoE5NQgecigDrSHkPtIYA",
	  "expires_in": 7200
	}
	*/
	var result = struct {
		AccessToken string  `json:"component_access_token"`
		ExpiresIn   int     `json:"expires_in"`
		Errcode     float64 `json:"errcode"`
		Errmsg      string  `json:"errmsg"`
	}{}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		err = fmt.Errorf("Unmarshal error %s", string(resp))
		return
	}

	if result.AccessToken == "" {
		err = fmt.Errorf("%s", string(resp))
		return
	}

	return result.AccessToken, result.ExpiresIn, nil
}

// GetComponentVerifyTicket 获取 component_verify_ticket
func GetComponentVerifyTicket(platform *Platform) (appTicket string, err error) {
	appTicket, err = platform.Cache.Fetch("component_verify_ticket :" + platform.Config.AppId)
	if appTicket != "" {
		return
	}

	return
}

// ReceiveComponentVerifyTicket 接收 component_verify_ticket
func ReceiveComponentVerifyTicket(platform *Platform, ticket string) (err error) {
	return platform.Cache.Save("component_verify_ticket :"+platform.Config.AppId, ticket, 0)
}
