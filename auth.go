package synology

const (
	authVersion  = 6
	authApiName  = "SYNO.API.Auth"
	loginMethod  = "login"
	logoutMethod = "logout"
)

type AuthResponse struct {
	Sid string
}

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
		"format":  "sid",
	}
	resp, err := s.makeRequest(login.Path, login.Name, loginMethod, authVersion, params)
	if err != nil {
		return err
	}

	result, err := unmarshal[AuthResponse](resp)
	s.sid = result.Sid

	return err
}

func (s *SynologyCore) Logout() error {
	login, err := s.Find(authApiName)
	if err != nil {
		return err
	}

	_, err = s.makeRequest(login.Path, login.Name, logoutMethod, authVersion, nil)

	return err
}
