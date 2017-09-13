package token

import "time"

type Token struct {
	Token        string
	RefreshToken string
	ExpireTime   time.Time
}
