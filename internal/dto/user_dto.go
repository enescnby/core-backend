package dto

type LookupResponse struct {
	UserID              string `json:"user_id"`
	EncryptionPublicKey string `json:"encryption_public_key"`
}
