package auth_user

import "context"

func (h *HTTPAuthService) ValidateTokenOrSession(ctx context.Context, token, sessionID, deviceID string) (*AuthResult, error) {
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
			UserID:   verifyJWTResponse.UserId,
			DeviceID: verifyJWTResponse.DeviceId,
			Valid:    verifyJWTResponse.Valid,
			Error:    verifyJWTResponse.Error,
		}
	} else if sessionID != "" {
		// validate session
		verifySessionResponse, err := h.authService.VerifySession(ctx, sessionID)
		if err != nil {
			return &AuthResult{
				Valid: false,
				Error: "http verify session error",
			}, err
		}

		result = &AuthResult{
			UserID:   verifySessionResponse.UserId,
			DeviceID: deviceID,
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

func (h *HTTPAuthService) RefreshSession(ctx context.Context, sessionID string) (*AuthResult, error) {
	refreshSessionResponse, err := h.authService.RefreshSession(ctx, sessionID)
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
