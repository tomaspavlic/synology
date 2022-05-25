package synology

import "errors"

const (
	authVersion  = 6
	authApiName  = "SYNO.API.Auth"
	loginMethod  = "login"
	logoutMethod = "logout"
)

// Login performs login with provided credentials. When successful session cookie is stored
// and used for all API requests made with SynologyCore client.
func (s *SynologyCore) Login(account, password string) error {
	login, err := s.Find(authApiName)
	if err != nil {
		return err
	}

	params := map[string]string{
		"account": account,
		"passwd":  password,
		"format":  "cookie",
	}
	resp, err := s.makeRequest(login.Path, login.Name, loginMethod, authVersion, params)
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

func (s *SynologyCore) Logout() error {
	login, err := s.Find(authApiName)
	if err != nil {
		return err
	}

	resp, err := s.makeRequest(login.Path, login.Name, logoutMethod, authVersion, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
