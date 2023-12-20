package models

type RefreshSession struct {
	UserID       string `redis:"userId"`
	RefreshToken string `redis:"refreshToken"`
	ExpiresAt    int64  `redis:"expiresAt"`
}
