package synology

const (
	authApiVersion      = 6
	authApiName         = "SYNO.API.Auth"
	authApiLoginMethod  = "login"
	authApiLogoutMethod = "logout"
)

type authResponse struct {
	Sid string
}

// Login performs login with provided credentials. When successful session cookie is stored
// and used for all API requests made with SynologyCore client.
func (s *SynologyCore) Login(account, password string) error {

	params := parameters{
		"account": account,
		"passwd":  password,
		"format":  "sid",
	}

	resp, err := s.makeRequest(
		authPath,
		authApiName,
		authApiLoginMethod,
		authApiVersion,
		params)

	if err != nil {
		return err
	}

	result, err := unmarshal[authResponse](resp)
	s.sid = result.Sid

	return err
}

func (s *SynologyCore) Logout() error {
	_, err := s.makeRequest(
		authPath,
		authApiName,
		authApiLogoutMethod,
		authApiVersion,
		nil,
	)

	return err
}
