package store

type (
	UserTokens struct {
		AccessToken  string
		RefreshToken string
		IDToken      string
	}
	CodeFlow struct {
		Verifier string
	}

	UserSession struct {
		ID       *string
		Tokens   *UserTokens
		CodeFlow *CodeFlow
	}

	Store interface {
		Get(sessionID string) *UserSession
		Start() (string, *UserSession)
		Update(sessionId string, session *UserSession) *UserSession
		Destroy(sessionId string)
	}
)

func NewUserSession(sessionID string) *UserSession {
	return &UserSession{
		ID: &sessionID,
	}
}

func NewUserTokens(accessToken string, refreshToken string, idToken string) *UserTokens {
	return &UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IDToken:      idToken,
	}
}

func NewCodeVerifier(verifier string) *CodeFlow {
	return &CodeFlow{
		Verifier: verifier,
	}
}
