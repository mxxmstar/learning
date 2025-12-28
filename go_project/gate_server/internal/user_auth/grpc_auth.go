package auth_user

import "context"

func (g *GRPCAuthService) ValidateTokenOrSession(ctx context.Context, token, sessionID, deviceID string) (*AuthResult, error) {
	var result *AuthResult

	if token != "" {
		verifyJWTResponse, err := g.authService.VerifyJWT(ctx, token)
		if err != nil {
			return &AuthResult{
				Valid: false,
				Error: "grpc verify jwt error",
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
		verifySessionResponse, err := g.authService.VerifySession(ctx, sessionID)
		if err != nil {
			return &AuthResult{
				Valid: false,
				Error: "grpc verify session error",
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

func (g *GRPCAuthService) RefreshSession(ctx context.Context, sessionID string) (*AuthResult, error) {
	refreshSessionResponse, err := g.authService.RefreshSession(ctx, sessionID)
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
