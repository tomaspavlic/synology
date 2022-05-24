package main

import "errors"

const (
	apiAuth     = "SYNO.API.Auth"
	loginMethod = "login"
)

func (s *SynologyCore) Login(account, password string) error {
	login, err := s.Find(apiAuth)
	if err != nil {
		return err
	}

	params := map[string]string{
		"account": account,
		"passwd":  password,
		"format":  "cookie",
	}
	resp, err := s.makeRequest(login.Path, login.Name, loginMethod, 1, params)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "id" {
			s.authCookie = cookie
			return nil
		}
	}

	return errors.New("auth cookie not found in auth response")
}
