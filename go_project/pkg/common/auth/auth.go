package common_auth

type VerifySessionRequest struct {
	SessionID string `json:"sessionId"`
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
	SessionID string `json:"sessionId"`
}

type RefreshSessionResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
