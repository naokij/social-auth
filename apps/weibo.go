// Copyright 2014 beego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
//
// Maintain by https://github.com/slene

package apps

import (
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"github.com/naokij/social-auth"
	"io/ioutil"
	"net/url"
)

type Weibo struct {
	BaseProvider
}

type WeiboUserInfo struct {
	Id          string `json:"id"`
	ScreenName  string `json:"screen_name"`
	AvatarLarge string `json:"avatar_large"`
}

func (u *WeiboUserInfo) GetLogin() string {
	return u.ScreenName
}

func (u *WeiboUserInfo) GetId() string {
	return u.Id
}

func (u *WeiboUserInfo) GetAvatarUrl() string {
	return u.AvatarLarge
}

func (u *WeiboUserInfo) GetEmail() string {
	return ""
}

func (p *Weibo) GetType() social.SocialType {
	return social.SocialWeibo
}

func (p *Weibo) GetName() string {
	return "Weibo"
}

func (p *Weibo) GetPath() string {
	return "weibo"
}

func (p *Weibo) GetIndentify(tok *social.Token) (string, error) {
	return tok.GetExtra("uid"), nil
}

func (p *Weibo) GetUserInfo(identity string, tok *social.Token) (userInfo social.UserInfo, err error) {
	uri := "https://api.weibo.com/2/users/show.json?"
	data := url.Values{}
	data.Add("uid", identity)
	data.Add("access_token", tok.AccessToken)
	req := httplib.Get(uri + data.Encode())
	resp, err := req.Response()
	if err != nil {
		return userInfo, err
	}
	defer resp.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return userInfo, err
	}
	weiboUserInfo := WeiboUserInfo{}
	userInfo = social.UserInfo(&weiboUserInfo)
	err = json.Unmarshal(body, &userInfo)
	return userInfo, err
	return
}

var _ social.Provider = new(Weibo)

func NewWeibo(clientId, secret string) *Weibo {
	p := new(Weibo)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "email"
	p.AuthURL = "https://api.weibo.com/oauth2/authorize"
	p.TokenURL = "https://api.weibo.com/oauth2/access_token"
	p.RedirectURL = social.DefaultAppUrl + "login/weibo/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}
