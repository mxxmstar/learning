package auth_user

import "context"

func (h *HTTPAuthService) ValidateTokenOrSession(ctx context.Context, token, sessionId, deviceId string) (*AuthResult, error) {
	var result *AuthResult

	if token != "" {
		verifyJWTResponse, err := h.authService.VerifyJWT(ctx, token)
		if err != nil {
			return &AuthResult{
				Valid: false,
				Error: "http verify jwt error",
			}, err
		}

		result = &AuthResult{
			UserId:   verifyJWTResponse.UserId,
			DeviceId: verifyJWTResponse.DeviceId,
			Valid:    verifyJWTResponse.Valid,
			Error:    verifyJWTResponse.Error,
		}
	} else if sessionId != "" {
		// validate session
		verifySessionResponse, err := h.authService.VerifySession(ctx, sessionId)
		if err != nil {
			return &AuthResult{
				Valid: false,
				Error: "http verify session error",
			}, err
		}

		result = &AuthResult{
			UserId:   verifySessionResponse.UserId,
			DeviceId: deviceId,
			Valid:    verifySessionResponse.Valid,
			Error:    verifySessionResponse.Error,
		}
	} else {
		return &AuthResult{
			Valid: false,
			Error: "no token and session id provided",
		}, nil
	}

	return result, nil
}

func (h *HTTPAuthService) RefreshSession(ctx context.Context, sessionId string) (*AuthResult, error) {
	refreshSessionResponse, err := h.authService.RefreshSession(ctx, sessionId)
	if err != nil {
		return &AuthResult{
			Valid: false,
			Error: "grpc refresh session error",
		}, err
	}

	return &AuthResult{
		Valid: refreshSessionResponse.Success,
		Error: refreshSessionResponse.Error,
	}, nil
}
