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
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"time"

	"github.com/naokij/social-auth"
)

type Github struct {
	BaseProvider
}

type GithubUserInfo struct {
	Login             string    `json:"login"`
	Id                int64     `json:"id"`
	AvatarUrl         string    `json:"avatar_url"`
	GravatarId        string    `json:"gravatar_id"`
	Url               string    `json:"url"`
	HtmlUrl           string    `json:"html_url"`
	FollowersUrl      string    `json:"followers_url"`
	FollowingUrl      string    `json:"following_url"`
	GistsUrl          string    `json:"gists_url"`
	StarredUrl        string    `json:"starred_url"`
	SubscriptionsUrl  string    `json:"subscriptions_url"`
	OrganizationsUrl  string    `json:"organizations_url"`
	ReposUrl          string    `json:"repos_url"`
	EventsUrl         string    `json:"events_url"`
	ReceivedEventsUrl string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	PublicRepos       int64     `json:"public_repos"`
	PublicGists       int64     `json:"public_gists"`
	Followers         int64     `json:"followers"`
	Following         int64     `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (u *GithubUserInfo) GetLogin() string {
	return u.Login
}

func (u *GithubUserInfo) GetId() string {
	return fmt.Sprintf("%d", u.Id)
}

func (u *GithubUserInfo) GetAvatarUrl() string {
	return u.AvatarUrl
}

func (u *GithubUserInfo) GetEmail() string {
	return u.Email
}

func (p *Github) GetType() social.SocialType {
	return social.SocialGithub
}

func (p *Github) GetName() string {
	return "Github"
}

func (p *Github) GetPath() string {
	return "github"
}

func (p *Github) GetIndentify(tok *social.Token) (string, error) {
	vals := make(map[string]interface{})

	uri := "https://api.github.com/user"
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)
	req.Header("Authorization", "Bearer "+tok.AccessToken)

	resp, err := req.Response()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()

	if err := decoder.Decode(&vals); err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%v", vals["message"])
	}

	if vals["id"] == nil {
		return "", nil
	}

	return fmt.Sprint(vals["id"]), nil
}

func (p *Github) GetUserInfo(identity string, tok *social.Token) (userInfo social.UserInfo, err error) {
	uri := "https://api.github.com/user"
	req := httplib.Get(uri)
	req.SetTransport(social.DefaultTransport)
	req.Header("Authorization", "Bearer "+tok.AccessToken)

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
	githubUserInfo := GithubUserInfo{}
	userInfo = social.UserInfo(&githubUserInfo)
	err = json.Unmarshal(body, &userInfo)
	return userInfo, err
}

var _ social.Provider = new(Github)

func NewGithub(clientId, secret string) *Github {
	p := new(Github)
	p.App = p
	p.ClientId = clientId
	p.ClientSecret = secret
	p.Scope = "user,public_repo"
	p.AuthURL = "https://github.com/login/oauth/authorize"
	p.TokenURL = "https://github.com/login/oauth/access_token"
	p.RedirectURL = social.DefaultAppUrl + "login/github/access"
	p.AccessType = "offline"
	p.ApprovalPrompt = "auto"
	return p
}
