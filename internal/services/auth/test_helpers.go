package auth

import (
	"auth_service/internal/domain/models"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"log"
	"testing"
	"time"
)

// getPublicKey creates an instance of rsa.PublicKey via const public key
func getPublicKey() *rsa.PublicKey {
	const pubKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy3FbmEBb1gjTFi2seOzj
kTJXvxQFeSCLUs+kczKRoKEWkU9C+07CICiG2meM4+scR3I9A7KbmrP976dJNMZc
NXxjAULLOkc9ZOJkHESs6lqKkQDfFihDB9i8LjRUSw6f14sawvTS6aGrYrKRpL3X
1oPtauUduHvyWfRvY+hz9SgY8WwjP5SXXy+b7sYStQM26CzSGIL53d3C33IqrHJh
ZSblLR5X2BpaDHCHMQNiBQhQ6vMMLsec96qtsfU8bcSUOmkplCf7LmDQZplS4K4S
KNv2mdTLj8GOWm25amkG1KYTPrHam65e+/BPmSBJKzr3qIu6oQTyspuBANZRksP4
AQIDAQAB
-----END PUBLIC KEY-----`

	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		log.Fatal("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse DER encoded public key: %v", err)
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub
	default:
		log.Fatal("unknown type of public key")
		return nil
	}
}

// getPrivateKey creates an instance of rsa.PublicKey via const public key
func getPrivateKey() *rsa.PrivateKey {
	const privKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDLcVuYQFvWCNMW
Lax47OORMle/FAV5IItSz6RzMpGgoRaRT0L7TsIgKIbaZ4zj6xxHcj0Dspuas/3v
p0k0xlw1fGMBQss6Rz1k4mQcRKzqWoqRAN8WKEMH2LwuNFRLDp/XixrC9NLpoati
spGkvdfWg+1q5R24e/JZ9G9j6HP1KBjxbCM/lJdfL5vuxhK1AzboLNIYgvnd3cLf
ciqscmFlJuUtHlfYGloMcIcxA2IFCFDq8wwux5z3qq2x9TxtxJQ6aSmUJ/suYNBm
mVLgrhIo2/aZ1MuPwY5abblqaQbUphM+sdqbrl778E+ZIEkrOveoi7qhBPKym4EA
1lGSw/gBAgMBAAECggEAEYcyQ5Nj9jRrb3E+92o9jyB9w+ZHNHDz4A4o1o4IUwUe
CyX/mORnwFqNh9q9HlP+6z7x99QXKQSR/+hzW7stbaRKaWzmUp6ZsQv6YR8foM9t
OeIP1npgkBgB2p9ClfbsGaeQBjUOyPdXa1kESGPc6UwTlA6qeV8gedSgFkUMXJZQ
2qdSGwJTNoUTt+JsXRaVeefjOsoTbEtRmKT/N1as9Pigm/k7vrf8ObB6IOC+RjKQ
W8ph9v/IETeJmphvCa0tz2SsjsZXixE/8ROKl30/G7fwgYQMniBX5KUnUDMGDmNT
RF9ARIO89v+TIa/+Gblnmv2r0Br3XtP3wjal1XtDgwKBgQD46FLCmJPe+ECdqNxm
5yUrmwR+e5UD/nIU/jJpOImeO5J68gy31TOu8nelD0TnXjL6ZocgVZfUhJ8bJhNA
B5uKs1iLlicTGjcV8TDwPz/tstaFwpv5xXeufX+HdHL4DCR4H0kyqs2Eki/MgY+B
rrb3OP64Btln9pciPqacFVSz2wKBgQDRPWNs1VGMLpxgf82BTKguDaU55gAQMn3A
VYN/J1arpeo2kE+jJE35HrNpUCcz2//RI9dQo8R88HewSDESXDizuPW76uAitGVw
Fv4VSGqxPI2rtsHwC5OpL6Wt4C8JKlP3ZhzBEZbc349++iK7NnIfp9M/Udf2d2fW
X8jtC614UwKBgQCR146mdsAt5Uf3GPoLUWR2KF55ve+SZ4RwyIDBJl98V2t8nlbV
YBboaymvjULSTl+QWILUb1KHMy4GukiNO+fnXS6Em3ZJuxKLyMbj/it3G1KXDXBW
6V024FHZDGJQ9MxpletNxMam5wEa0s9DLRwHv12AdoLsZ5AmgI8e5WC8AQKBgCnt
LxsDs49vV45OjZM3FQwFV/I+EA0u0NvVRsAX1doXKNM+H3cFM0qTyEd19CUEFEKH
+AEEj76pQJJmJM8VA5efnD3HVpodo6XONaN2G0rgY1LhEANVjoT8MyqgHUys2p5c
K67UI3KmqU48OPFie4O0UTQC5k9QpdHi55P4Hw5ZAoGAe5UEqjG1hs+75v07Vomg
ycelIcJavsscmDji/cmvRAQ1I5ymn+7SfYE20K9NWhkDeiP6eu86cHI78XrQdh+q
ZVljiv3QbnCAfXdyeCclGTJrDH8s5MeBew8TtgXO+XnENgkwv+dL2WLRJz95ozFW
F7fYt14Djv6XwFFVuEKSq2s=
-----END PRIVATE KEY-----`

	block, _ := pem.Decode([]byte(privKey))
	if block == nil {
		log.Fatal("failed to parse PEM block containing the private key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse DER encoded private key: %v", err)
	}

	switch priv := priv.(type) {
	case *rsa.PrivateKey:
		return priv
	default:
		log.Fatal("unknown type of private key")
		return nil
	}
}

// MockRefreshSessionsRepository is a mock implementation of the RefreshSessionsRepository interface
type MockRefreshSessionsRepository struct {
	mock.Mock
}

func (m *MockRefreshSessionsRepository) GetRefreshSession(ctx context.Context, userId string) (*models.RefreshSession, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*models.RefreshSession), args.Error(1)
}

func (m *MockRefreshSessionsRepository) InsertRefreshSession(ctx context.Context, session *models.RefreshSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockRefreshSessionsRepository) DeleteRefreshSession(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockRefreshSessionsRepository) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

// newTestJwtAuthService creates a instance JwtAuthService with test-env values.
// Uses mocked redisRepository and zaptest-version of logger.
func newTestJwtAuthService(t *testing.T) *JwtAuthService {
	// Creating mocks
	mockRedisRepo := new(MockRefreshSessionsRepository)
	testLogger := zaptest.NewLogger(t).Sugar()

	// Setting expected behavior of mocks
	mockRedisRepo.On("InsertRefreshSession", mock.Anything, mock.AnythingOfType("*models.RefreshSession")).Return(nil)

	// Creating an instance of JwtAuthService
	authService := &JwtAuthService{
		privateKey:      getPrivateKey(),
		publicKey:       getPublicKey(),
		signingMethod:   jwt.SigningMethodRS256,
		redisRepository: mockRedisRepo,
		log:             testLogger,
		refreshTokenTTL: time.Hour * 24,
		tokenTTL:        time.Minute * 15,
	}

	return authService
}
