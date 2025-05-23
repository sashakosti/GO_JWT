package token

type Validator interface {
	Validate(token string) error
}

type TokenValidator struct {
	access  *AccessTokenManager
	refresh *RefreshTokenManager
}

func NewTokenValidator(access *AccessTokenManager, refresh *RefreshTokenManager) *TokenValidator {
	return &TokenValidator{
		access:  access,
		refresh: refresh,
	}
}

func (v *TokenValidator) ValidateAccess(token string) (string, error) {
	return v.access.Validate(token)
}

func (v *TokenValidator) ValidateRefresh(userID, token string) error {
	return v.refresh.Validate(userID, token)
}
