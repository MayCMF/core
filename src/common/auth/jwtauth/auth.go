package jwtauth

import (
	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/MayCMF/core/src/common/auth"
)

const defaultKey = "GINADMIN"

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       7200,
	signingMethod: jwt.SigningMethodHS512,
	signingKey:    []byte(defaultKey),
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(defaultKey), nil
	},
}

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       int
	tokenType     string
}

// Option - Defining parameter items
type Option func(*options)

// SetSigningMethod - Set signature method
func SetSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// SetSigningKey - Set signature key
func SetSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

// SetKeyfunc - Set the callback function of the verification key
func SetKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfunc = keyFunc
	}
}

// SetExpired - Set the token expiration time (in seconds, default 7200)
func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

// New - Create an authentication instance
func New(store Storer, opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &JWTAuth{
		opts:  &o,
		store: store,
	}
}

// JWTAuth - jwt Authenticate
type JWTAuth struct {
	opts  *options
	store Storer
}

// GenerateToken - Generate token
func (a *JWTAuth) GenerateToken(ctx context.Context, userUUID string) (auth.TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(a.opts.expired) * time.Second).Unix()

	token := jwt.NewWithClaims(a.opts.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   userUUID,
	})

	tokenString, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return nil, err
	}

	tokenInfo := &tokenInfo{
		ExpiresAt:   expiresAt,
		TokenType:   a.opts.tokenType,
		AccessToken: tokenString,
	}
	return tokenInfo, nil
}

// Parsing token
func (a *JWTAuth) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, a.opts.keyfunc)
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

func (a *JWTAuth) callStore(fn func(Storer) error) error {
	if store := a.store; store != nil {
		return fn(store)
	}
	return nil
}

// DestroyToken - Destroy token
func (a *JWTAuth) DestroyToken(ctx context.Context, tokenString string) error {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return err
	}

	// If the storage is set, the unexpired token is placed
	return a.callStore(func(store Storer) error {
		expired := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
		return store.Set(ctx, tokenString, expired)
	})
}

// ParseUserUUID - Resolve user UUID
func (a *JWTAuth) ParseUserUUID(ctx context.Context, tokenString string) (string, error) {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return "", err
	}

	err = a.callStore(func(store Storer) error {
		exists, err := store.Check(ctx, tokenString)
		if err != nil {
			return err
		} else if exists {
			return auth.ErrInvalidToken
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

// Release - Release resources
func (a *JWTAuth) Release() error {
	return a.callStore(func(store Storer) error {
		return store.Close()
	})
}
