package common_auth

type VerifySessionRequest struct {
	SessionId string `json:"sessionId"`
}

type VerifySessionResponse struct {
	Valid  bool   `json:"valid"`
	UserId uint64 `json:"userId"`
	Error  string `json:"error,omitempty"`
}

type VerifyJWTRequest struct {
	JWTToken string `json:"jwtToken"`
}

type VerifyJWTResponse struct {
	Valid    bool   `json:"valid"`
	UserId   uint64 `json:"userId,omitempty"`
	DeviceId string `json:"deviceId,omitempty"`
	Error    string `json:"error,omitempty"`
}

type RefreshSessionRequest struct {
	SessionId string `json:"sessionId"`
}

type RefreshSessionResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type LoginByEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceId string `json:"deviceId"`
}

type LoginByEmailResponse struct {
	SessionId string `json:"sessionId"`
	JWTToken  string `json:"jwtToken"`
	UserId    uint64 `json:"userId"`
	Error     string `json:"error,omitempty"`
}

type SignUpRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type SignUpResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
