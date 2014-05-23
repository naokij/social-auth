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
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/astaxie/beego/httplib"

	"github.com/naokij/social-auth"
)

type QQ struct {
	BaseProvider
}

type QQUserInfo struct {
	Ret             string `json:"ret"`
	Msg             string `json:"msg"`
	Nickname        string `json:"nickname"`
	Figureurl       string `json:"figureurl"`
	Figureurl1      string `json:"figureurl_1"`
	Figureurl2      string `json:"figureurl_2"`
	FigureurlQQ1    string `json:"figureurl_qq_1"`
	FigureurlQQ2    string `json:"figureurl_qq_2"`
	Gender          string `json:"gender"`
	IsYellowVip     int    `json:"is_yellow_vip"`
	Vip             int    `json:"vip"`
	YellowVipLevel  int    `json:"yellow_vip_level"`
	Level           int    `json:"level"`
	IsYellowYearVip int    `json:"is_yellow_year_vip"`
}

func (u *QQUserInfo) GetLogin() string {
	return u.Nickname
}

func (u *QQUserInfo) GetId() string {
	return ""
}

func (u *QQUserInfo) GetAvatarUrl() string {
	return u.FigureurlQQ2
}

func (u *QQUserInfo) GetEmail() string {
	return ""
}

func (p *QQ) GetType() social.SocialType {
	return social.SocialQQ
}

func (p *QQ) GetName() string {
	return "QQ"
}

func (p *QQ) GetPath() string {
	return "qq"
}

func (p *QQ) GetIndentify(tok *social.Token) (string, error) {
	uri := "https://graph.z.qq.com/moc2/me?access_token=" + url.QueryEscape(tok.AccessToken)
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)

	body, err := req.String()
	if err != nil {
		return "", err
	}

	vals, err := url.ParseQuery(body)
	if err != nil {
		return "", err
	}

	if vals.Get("code") != "" {
		return "", fmt.Errorf("code: %s, msg: %s", vals.Get("code"), vals.Get("msg"))
	}

	return vals.Get("openid"), nil
}

func (p *QQ) GetUserInfo(identity string, tok *social.Token) (userInfo social.UserInfo, err error) {
	uri := "https://api.github.com/user?"
	data := url.Values{}
	data.Add("access_token", tok.AccessToken)
	data.Add("oauth_consumer_key", p.ClientSecret)
	data.Add("openid", identity)
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
	qqUserInfo := QQUserInfo{}
	userInfo = social.UserInfo(&qqUserInfo)
	err = json.Unmarshal(body, &userInfo)
	return userInfo, err
}

var _ social.Provider = new(QQ)

func NewQQ(clientId, secret string) *QQ {
	p := new(QQ)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "get_user_info"
	p.AuthURL = "https://graph.qq.com/oauth2.0/authorize"
	p.TokenURL = "https://graph.qq.com/oauth2.0/token"
	p.RedirectURL = social.DefaultAppUrl + "login/qq/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}
