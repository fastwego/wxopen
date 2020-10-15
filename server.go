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
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fastwego/offiaccount/type/type_message"
	"github.com/fastwego/offiaccount/util"
	"github.com/fastwego/wxopen/type/type_platform"
)

/*
响应微信请求 或 推送消息/事件 的服务器
*/
type Server struct {
	Ctx *Platform
}

// ParseXML 解析微信推送过来的消息/事件
func (s *Server) ParseXML(body []byte) (m interface{}, err error) {

	if s.Ctx.Logger != nil {
		s.Ctx.Logger.Println(string(body))
	}

	// 是否加密消息
	encryptMsg := type_message.EncryptMessage{}
	err = xml.Unmarshal(body, &encryptMsg)
	if err != nil {
		return
	}

	// 需要解密
	if encryptMsg.Encrypt != "" {
		var xmlMsg []byte
		_, xmlMsg, _, err = util.AESDecryptMsg(encryptMsg.Encrypt, s.Ctx.Config.AesKey)
		if err != nil {
			return
		}
		body = xmlMsg

		if s.Ctx.Logger != nil {
			s.Ctx.Logger.Println("AESDecryptMsg ", string(body))
		}
	}

	event := type_platform.Event{}
	err = xml.Unmarshal(body, &event)
	//fmt.Println(message)
	if err != nil {
		return
	}

	switch event.InfoType {

	case type_platform.EventTypeComponentVerifyTicket:
		msg := type_platform.EventComponentVerifyTicket{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case type_platform.EventTypeAuthorized:
		msg := type_platform.EventAuthorized{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case type_platform.EventTypeUnauthorized:
		msg := type_platform.EventUnauthorized{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case type_platform.EventTypeUpdateAuthorized:
		msg := type_platform.EventUpdateAuthorized{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil

	}
	return
}

// Response 响应微信消息 (自动判断是否要加密)
func (s *Server) Response(writer http.ResponseWriter, request *http.Request, reply interface{}) (err error) {

	output := []byte("success") // 默认回复
	if reply != nil {
		output, err = xml.Marshal(reply)
		if err != nil {
			return
		}

		// 加密
		if request.URL.Query().Get("encrypt_type") == "aes" {
			message := s.encryptReplyMessage(output)
			output, err = xml.Marshal(message)
			if err != nil {
				return
			}
		}
	}

	_, err = writer.Write(output)

	if s.Ctx.Logger != nil {
		s.Ctx.Logger.Println("Response: ", string(output))
	}

	return
}

// encryptReplyMessage 加密回复消息
func (s *Server) encryptReplyMessage(rawXmlMsg []byte) (replyEncryptMessage type_message.ReplyEncryptMessage) {
	cipherText := util.AESEncryptMsg([]byte(util.GetRandString(16)), rawXmlMsg, s.Ctx.Config.AppId, s.Ctx.Config.AesKey)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := util.GetRandString(6)

	strs := []string{
		timestamp,
		nonce,
		s.Ctx.Config.Token,
		cipherText,
	}
	sort.Strings(strs)
	h := sha1.New()
	_, _ = io.WriteString(h, strings.Join(strs, ""))
	signature := fmt.Sprintf("%x", h.Sum(nil))

	return type_message.ReplyEncryptMessage{
		Encrypt:      cipherText,
		MsgSignature: signature,
		TimeStamp:    timestamp,
		Nonce:        nonce,
	}
}
