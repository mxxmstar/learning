package auth_user

import "context"

func (g *GRPCAuthService) ValidateTokenOrSession(ctx context.Context, token, sessionId, deviceId string) (*AuthResult, error) {
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
			UserId:   verifyJWTResponse.UserId,
			DeviceId: verifyJWTResponse.DeviceId,
			Valid:    verifyJWTResponse.Valid,
			Error:    verifyJWTResponse.Error,
		}
	} else if sessionId != "" {
		// validate session
		verifySessionResponse, err := g.authService.VerifySession(ctx, sessionId)
		if err != nil {
			return &AuthResult{
				Valid: false,
				Error: "grpc verify session error",
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

func (g *GRPCAuthService) RefreshSession(ctx context.Context, sessionId string) (*AuthResult, error) {
	refreshSessionResponse, err := g.authService.RefreshSession(ctx, sessionId)
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

func (g *GRPCAuthService) SignUp(ctx context.Context, username, email, password, confirmPassword string) (*AuthResult, error) {
	signUpResponse, err := g.authService.SignUp(ctx, username, email, password, confirmPassword)
	if err != nil {
		return &AuthResult{
			Valid: false,
			Error: "grpc sign up error",
		}, err
	}

	return &AuthResult{
		Valid: signUpResponse.Success,
		Error: signUpResponse.Error,
	}, nil
}
