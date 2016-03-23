/**
  Steam Library For Go
  Copyright (C) 2016 Ahmed Samy <f.fallen45@gmail.com>

  This library is free software; you can redistribute it and/or
  modify it under the terms of the GNU Lesser General Public
  License as published by the Free Software Foundation; either
  version 2.1 of the License, or (at your option) any later version.

  This library is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
  Lesser General Public License for more details.

  You should have received a copy of the GNU Lesser General Public
  License along with this library; if not, write to the Free Software
  Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/
package steam

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	apiKeyURL         = "https://steamcommunity.com/dev/apikey"
	apiKeyRegisterURL = "https://steamcommunity.com/dev/registerkey"
	apiKeyRevokeURL   = "https://steamcommunity.com/dev/revokekey"

	accessDeniedPattern = "<h2>Access Denied</h2>"
)

var (
	keyRegExp = regexp.MustCompile("<p>Key: ([0-9A-F]+)</p>")

	ErrCannotRegisterKey = errors.New("unable to register API key")
	ErrCannotRevokeKey   = errors.New("unable to revoke API key")
	ErrAccessDenied      = errors.New("access is denied")
	ErrKeyNotFound       = errors.New("key not found")
)

func (community *Community) RegisterWebAPIKey(domain string) error {
	values := url.Values{
		"domain":       {domain},
		"agreeToTerms": {"agreed"},
		"sessionid":    {community.sessionID},
		"Submit":       {"Register"},
	}

	req, err := http.NewRequest(http.MethodPost, apiKeyRegisterURL, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	resp, err := community.client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrCannotRegisterKey
	}

	return nil
}

func (community *Community) GetWebAPIKey() (string, error) {
	req, err := http.NewRequest(http.MethodGet, apiKeyURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := community.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if m, err := regexp.Match(accessDeniedPattern, body); err != nil {
		return "", err
	} else if m {
		return "", ErrAccessDenied
	}

	submatch := keyRegExp.FindStringSubmatch(string(body))
	if len(submatch) <= 1 {
		return "", ErrKeyNotFound
	}

	community.apiKey = submatch[1]
	return submatch[1], nil
}

func (community *Community) RevokeWebAPIKey() error {
	values := url.Values{
		"revoke":    {"Revoke My Steam Web API Key"},
		"sessionid": {community.sessionID},
	}

	req, err := http.NewRequest(http.MethodPost, apiKeyRevokeURL, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	resp, err := community.client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrCannotRevokeKey
	}

	return nil
}
