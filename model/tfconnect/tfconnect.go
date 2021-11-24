// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package oauthtfconnect

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"
)

type TFConnectProvider struct {
}

type TFConnectUser struct {
	Username string `json:"username"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func init() {
	provider := &TFConnectProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_TFCONNECT, provider)
}

func userFromTFConnectUser(tfu *TFConnectUser) *model.User {
	user := &model.User{}
	username := tfu.Username
	user.Username = model.CleanUsername(username)
	user.FirstName = user.Username
	user.Email = tfu.Email
	user.Email = strings.ToLower(user.Email)
	user.AuthData = &user.Email
	user.AuthService = model.USER_AUTH_SERVICE_TFCONNECT

	return user
}

func tfConnectUserFromJson(data io.Reader) (*TFConnectUser, error) {
	decoder := json.NewDecoder(data)
	var tfu TFConnectUser
	err := decoder.Decode(&tfu)
	if err != nil {
		return nil, err
	}
	return &tfu, nil

}

func (tfu *TFConnectUser) ToJson() string {
	b, err := json.Marshal(tfu)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (tfu *TFConnectUser) IsValid() error {
	// if tfu.Id == 0 {
	// 	return errors.New("user id can't be 0")
	// }

	if tfu.Email == "" {
		return errors.New("user e-mail should not be empty")
	}

	return nil
}

func (m *TFConnectProvider) GetUserFromJson(data io.Reader, tokenUser *model.User) (*model.User, error) {
	tfu, err := tfConnectUserFromJson(data)
	if err != nil {
		return nil, err
	}
	if err = tfu.IsValid(); err != nil {
		return nil, err
	}

	return userFromTFConnectUser(tfu), nil
}

func (m *TFConnectProvider) GetSSOSettings(config *model.Config, service string) (*model.SSOSettings, error) {
	return &config.TFConnectSettings, nil
}

func (m *TFConnectProvider) GetUserFromIdToken(idToken string) (*model.User, error) {
	return nil, nil
}

func (m *TFConnectProvider) IsSameUser(dbUser, oauthUser *model.User) bool {
	return dbUser.AuthData == oauthUser.AuthData
}
