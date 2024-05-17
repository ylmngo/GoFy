package data

import (
	"strconv"
	"time"

	"github.com/pascaldekloe/jwt"
)

func CreateJWT(sub int, expires time.Time, issuer string, audiences []string, secret string) ([]byte, error) {
	var claims jwt.Claims

	claims.Subject = strconv.FormatInt(int64(sub), 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expires)
	claims.Audiences = audiences
	claims.Issuer = issuer

	token, err := claims.HMACSign(jwt.HS256, []byte(secret))
	if err != nil {
		return nil, err
	}

	return token, nil
}
