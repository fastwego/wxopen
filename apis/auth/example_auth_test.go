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

package auth_test

import (
	"fmt"

	"github.com/fastwego/wxopen"
	"github.com/fastwego/wxopen/apis/auth"
)

func ExampleCreatePreauthCode() {
	var ctx *wxopen.Platform

	payload := []byte("{}")
	resp, err := auth.CreatePreauthCode(ctx, payload)

	fmt.Println(resp, err)
}

func ExampleApiQueryAuth() {
	var ctx *wxopen.Platform

	payload := []byte("{}")
	resp, err := auth.ApiQueryAuth(ctx, payload)

	fmt.Println(resp, err)
}

func ExampleApiAuthorizerToken() {
	var ctx *wxopen.Platform

	payload := []byte("{}")
	resp, err := auth.ApiAuthorizerToken(ctx, payload)

	fmt.Println(resp, err)
}

func ExampleApiGetAuthorizerInfo() {
	var ctx *wxopen.Platform

	payload := []byte("{}")
	resp, err := auth.ApiGetAuthorizerInfo(ctx, payload)

	fmt.Println(resp, err)
}

func ExampleApiGetAuthorizerOption() {
	var ctx *wxopen.Platform

	payload := []byte("{}")
	resp, err := auth.ApiGetAuthorizerOption(ctx, payload)

	fmt.Println(resp, err)
}

func ExampleApiSetAuthorizerOption() {
	var ctx *wxopen.Platform

	payload := []byte("{}")
	resp, err := auth.ApiSetAuthorizerOption(ctx, payload)

	fmt.Println(resp, err)
}
