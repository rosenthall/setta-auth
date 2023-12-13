package models

type RefreshSession struct {
	ID           string `redis:"id"`
	RefreshToken string `redis:"refreshToken"`
	UserAgent    string `redis:"userAgent"`
	Fingerprint  string `redis:"fingerprint"`
	ExpiresAt    int64  `redis:"expiresAt"`
	CreatedAt    int64  `redis:"createdAt"`
}
