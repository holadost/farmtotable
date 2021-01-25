package aragorn

import (
	"context"
	firebase "firebase.google.com/go"
	firebase_auth "firebase.google.com/go/auth"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"google.golang.org/api/option"
	"time"
)

// Various auth related errors.
const (
	KNoIDToken = 1
	KInvalidIDToken
)

type AuthError struct {
	errorCode uint
	errorMsg  string
}

func NewAuthError(errorCode uint, errorMsg string) *AuthError {
	return &AuthError{
		errorCode: errorCode,
		errorMsg:  errorMsg,
	}
}

func (ae *AuthError) Error() string {
	return fmt.Sprintf(
		"AuthError(%s): %s",
		ae.ErrorCodeStr(), ae.errorMsg)
}

func (ae *AuthError) ErrorMsg() string {
	return ae.errorMsg
}

func (ae *AuthError) ErrorCode() uint {
	return ae.errorCode
}

func (ae *AuthError) ErrorCodeStr() string {
	if ae.errorCode == KNoIDToken {
		return "KNoIDToken"
	} else if ae.errorCode == KInvalidIDToken {
		return "KInvalidIDToken"
	} else {
		glog.Fatalf("Invalid error code: %d", ae.errorCode)
	}
	return ""
}

type authCacheEntry struct {
	token *firebase_auth.Token
}

type Auth struct {
	firebaseApp  *firebase.App
	firebaseAuth *firebase_auth.Client
	tokenCache   *ristretto.Cache
}

func NewAuth(credPath string) *Auth {
	auth := &Auth{}
	opt := option.WithCredentialsFile(*fbCredPath)
	var err error
	auth.firebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		glog.Fatalf("Unable to initialize firebase app due to err: %s", err.Error())
	}
	auth.firebaseAuth, err = auth.firebaseApp.Auth(context.Background())
	if err != nil {
		glog.Fatalf("Unable to initialize firebase auth client due to err: %s", err.Error())
	}
	auth.tokenCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 25000, // number of keys to track frequency of (10M).
		MaxCost:     2048,  // maximum cost of cache (1GB).
		BufferItems: 64,    // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return auth
}

func (auth *Auth) Authenticate(c *gin.Context) error {
	headers := c.Request.Header["Authorization"]
	if (headers == nil) || (len(headers) == 0) {
		return NewAuthError(KNoIDToken, "No authorization headers present")
	}
	idToken := headers[0]
	glog.V(1).Infof("ID Token: %s", idToken)
	entryIf, present := auth.tokenCache.Get(idToken)
	if present {
		entry := entryIf.(*authCacheEntry)
		c.Set("Token", entry.token)
		return nil
	}
	token, err := auth.firebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return NewAuthError(KInvalidIDToken, err.Error())
	}
	auth.tokenCache.SetWithTTL(
		idToken, &authCacheEntry{token: token}, 1,
		time.Duration(token.Expires-token.IssuedAt))
	c.Set("Token", token)
	return nil
}
