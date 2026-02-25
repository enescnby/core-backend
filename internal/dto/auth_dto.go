package dto

type RegisterRequest struct {
	PublicKey           string `json:"public_key"`
	EncryptedPrivateKey string `json:"encrypted_private_key"`
	Salt                string `json:"salt"`
	DeviceModel         string `json:"device_model"`
	FCMToken            string `json:"fcm_token"`
}

type RegisterResponse struct {
	CoreGuardID string `json:"core_guard_id"`
	Message     string `json:"message"`
}

type LoginInitRequest struct {
	CoreGuardID string `json:"core_guard_id"`
}

type LoginInitResponse struct {
	EncryptedPrivateKey string `json:"encrypted_private_key"`
	Salt                string `json:"salt"`
	Challenge           string `json:"challenge"`
}

type LoginVerifyRequest struct {
	CoreGuardID string `json:"core_guard_id"`
	Challenge   string `json:"challenge"`
	Signature   string `json:"signature"`
	DeviceModel string `json:"device_model"`
	FCMToken    string `json:"fcm_token"`
}

type LoginVerifyResponse struct {
	AccessToken string `json:"access_token"`
	Message     string `json:"message"`
}
