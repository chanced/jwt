package jwt_test

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/chanced/go-jwt/v4"
	"github.com/chanced/go-jwt/v4/test"
)

var errKeyFuncError error = fmt.Errorf("error loading key")

var errMap = map[error]string{
	jwt.ErrMalformedToken:              "ErrMalformedToken",
	jwt.ErrTokenContainsBearer:         "ErrTokenContainsBearer",
	jwt.ErrInvalidSigningMethod:        "ErrInvalidSigningMethod",
	jwt.ErrUnregisteredSigningMethod:   "ErrUnregisteredSigningMethod",
	jwt.ErrInvalidKey:                  "ErrInvalidKey",
	jwt.ErrInvalidKeyType:              "ErrInvalidKeyType",
	jwt.ErrHashUnavailable:             "ErrHashUnavailable",
	jwt.ErrTokenNotYetValid:            "ErrTokenNotYetValid",
	jwt.ErrTokenExpired:                "ErrTokenExpired",
	jwt.ErrTokenUsedBeforeIssued:       "ErrTokenUsedBeforeIssued",
	jwt.ErrNoneSignatureTypeDisallowed: "ErrNoneSignatureTypeDisallowed",
	jwt.ErrMissingKeyFunc:              "ErrMissingKeyFunc",
	jwt.ErrSignatureInvalid:            "ErrSignatureInvalid",
	jwt.ErrKeyFuncError:                "ErrKeyFuncError",
}

var (
	jwtTestDefaultKey *rsa.PublicKey
	defaultKeyFunc    jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return jwtTestDefaultKey, nil }
	emptyKeyFunc      jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return nil, nil }
	errorKeyFunc      jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return nil, errKeyFuncError }
	nilKeyFunc        jwt.Keyfunc = nil
)

func init() {
	jwtTestDefaultKey = test.LoadRSAPublicKeyFromDisk("test/sample_key.pub")
}

type Errors []error

func (errs Errors) Contains(error error) bool {
	for _, err := range errs {
		if errors.Is(err, error) {
			return true
		}
	}
	return false
}

var jwtTestData = []struct {
	name        string
	tokenString string
	keyfunc     jwt.Keyfunc
	claims      jwt.Claims
	valid       bool
	errors      Errors
	parser      *jwt.Parser
}{
	{
		"basic",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		true,
		nil,
		nil,
	},
	{
		"basic expired",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "exp": float64(time.Now().Unix() - 100)},
		false,
		Errors{jwt.ErrTokenExpired},
		nil,
	},
	{
		"basic nbf",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "nbf": float64(time.Now().Unix() + 100)},
		false,
		Errors{jwt.ErrTokenNotYetValid},
		nil,
	},
	{
		"expired and nbf",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "nbf": float64(time.Now().Unix() + 100), "exp": float64(time.Now().Unix() - 100)},
		false,
		Errors{jwt.ErrTokenNotYetValid, jwt.ErrTokenExpired},
		nil,
	},
	{
		"basic invalid",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		false,
		Errors{jwt.ErrSignatureInvalid},
		nil,
	},
	{
		"basic nokeyfunc",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		nilKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		false,
		Errors{jwt.ErrMissingKeyFunc},
		nil,
	},
	{
		"basic nokey",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		emptyKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		false,
		Errors{jwt.ErrInvalidKeyType},
		nil,
	},
	{
		"basic errorkey",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		errorKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		false,
		Errors{jwt.ErrKeyFuncError, errKeyFuncError},
		nil,
	},
	{
		"invalid signing method",
		"",
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		false,
		Errors{jwt.ErrInvalidSigningMethod},
		&jwt.Parser{ValidMethods: []string{"HS256"}},
	},
	{
		"valid signing method",
		"",
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		true,
		nil,
		&jwt.Parser{ValidMethods: []string{"RS256", "HS256"}},
	},
	{
		"JSON Number",
		"",
		defaultKeyFunc,
		jwt.MapClaims{"foo": json.Number("123.4")},
		true,
		nil,
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"Standard Claims",
		"",
		defaultKeyFunc,
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 10).Unix(),
		},
		true,
		nil,
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"JSON Number - basic expired",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "exp": json.Number(fmt.Sprintf("%v", time.Now().Unix()-100))},
		false,
		Errors{jwt.ErrTokenExpired},
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"JSON Number - basic nbf",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100))},
		false,
		Errors{jwt.ErrTokenNotYetValid},
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"JSON Number - expired and nbf",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100)), "exp": json.Number(fmt.Sprintf("%v", time.Now().Unix()-100))},
		false,
		Errors{jwt.ErrTokenNotYetValid, jwt.ErrTokenExpired},
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"SkipClaimsValidation during token parsing",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100))},
		true,
		nil,
		&jwt.Parser{UseJSONNumber: true, SkipClaimsValidation: true},
	},
	{
		"RFC7519 Claims",
		"",
		defaultKeyFunc,
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 10)),
		},
		true,
		nil,
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"RFC7519 Claims - single aud",
		"",
		defaultKeyFunc,
		&jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{"test"},
		},
		true,
		nil,
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"RFC7519 Claims - multiple aud",
		"",
		defaultKeyFunc,
		&jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{"test", "test"},
		},
		true,
		nil,
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"RFC7519 Claims - single aud with wrong type",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOjF9.8mAIDUfZNQT3TGm1QFIQp91OCpJpQpbB1-m9pA2mkHc", // { "aud": 1 }
		defaultKeyFunc,
		&jwt.RegisteredClaims{
			Audience: nil, // because of the unmarshal error, this will be empty
		},
		false,
		Errors{jwt.ErrMalformedToken},
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"RFC7519 Claims - multiple aud with wrong types",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsidGVzdCIsMV19.htEBUf7BVbfSmVoTFjXf3y6DLmDUuLy1vTJ14_EX7Ws", // { "aud": ["test", 1] }
		defaultKeyFunc,
		&jwt.RegisteredClaims{
			Audience: nil, // because of the unmarshal error, this will be empty
		},
		false,
		Errors{jwt.ErrMalformedToken},
		&jwt.Parser{UseJSONNumber: true},
	},
}

func TestParser_Parse(t *testing.T) {
	privateKey := test.LoadRSAPrivateKeyFromDisk("test/sample_key")

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		t.Run(data.name, func(t *testing.T) {
			// If the token string is blank, use helper function to generate string
			if data.tokenString == "" {
				data.tokenString = test.MakeSampleToken(data.claims, privateKey)
			}

			// Parse the token
			var token *jwt.Token
			var err error
			var parser = data.parser
			if parser == nil {
				parser = new(jwt.Parser)
			}
			// Figure out correct claims type
			switch data.claims.(type) {
			case jwt.MapClaims:
				token, err = parser.ParseWithClaims(data.tokenString, jwt.MapClaims{}, data.keyfunc)
			case *jwt.StandardClaims:
				token, err = parser.ParseWithClaims(data.tokenString, &jwt.StandardClaims{}, data.keyfunc)
			case *jwt.RegisteredClaims:
				token, err = parser.ParseWithClaims(data.tokenString, &jwt.RegisteredClaims{}, data.keyfunc)
			}
			if token == nil {
				panic("token is nil")
			}
			// Verify result matches expectation
			if !reflect.DeepEqual(data.claims, token.Claims) {
				t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
			}

			if data.valid && err != nil {
				t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
			}

			if !data.valid && err == nil {
				t.Errorf("[%v] Invalid token passed validation", data.name)
			}

			if (err == nil && !token.Valid) || (err != nil && token.Valid) {
				t.Errorf("[%v] Inconsistent behavior between returned error and token.Valid", data.name)
			}

			if len(data.errors) != 0 {
				if err == nil {
					t.Errorf("[%v] Expecting error.  Didn't get one.", data.name)
				} else {
					for _, expectedError := range data.errors {
						if !errors.Is(err, expectedError) {
							t.Errorf(`[%v] Expected "%v", received: %v`, data.name, errMap[expectedError], err)
						}
					}
				}
			}
			if data.valid && token.Signature == "" {
				t.Errorf("[%v] Signature is left unpopulated after parsing", data.name)
			}
		})
	}
}

func TestParser_ParseUnverified(t *testing.T) {
	privateKey := test.LoadRSAPrivateKeyFromDisk("test/sample_key")

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		// Skip test data, that intentionally contains malformed tokens, as they would lead to an error
		if data.errors.Contains(jwt.ErrMalformedToken) {
			continue
		}

		t.Run(data.name, func(t *testing.T) {
			// If the token string is blank, use helper function to generate string
			if data.tokenString == "" {
				data.tokenString = test.MakeSampleToken(data.claims, privateKey)
			}

			// Parse the token
			var token *jwt.Token
			var err error
			var parser = data.parser
			if parser == nil {
				parser = new(jwt.Parser)
			}
			// Figure out correct claims type
			switch data.claims.(type) {
			case jwt.MapClaims:
				token, _, err = parser.ParseUnverified(data.tokenString, jwt.MapClaims{})
			case *jwt.StandardClaims:
				token, _, err = parser.ParseUnverified(data.tokenString, &jwt.StandardClaims{})
			case *jwt.RegisteredClaims:
				token, _, err = parser.ParseUnverified(data.tokenString, &jwt.RegisteredClaims{})
			}

			if err != nil {
				t.Errorf("[%v] Invalid token", data.name)
			}

			// Verify result matches expectation
			if !reflect.DeepEqual(data.claims, token.Claims) {
				t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
			}

			if data.valid && err != nil {
				t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
			}
		})
	}
}

func BenchmarkParseUnverified(b *testing.B) {
	privateKey := test.LoadRSAPrivateKeyFromDisk("test/sample_key")

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		// If the token string is blank, use helper function to generate string
		if data.tokenString == "" {
			data.tokenString = test.MakeSampleToken(data.claims, privateKey)
		}

		// Parse the token
		var parser = data.parser
		if parser == nil {
			parser = new(jwt.Parser)
		}
		// Figure out correct claims type
		switch data.claims.(type) {
		case jwt.MapClaims:
			b.Run("map_claims", func(b *testing.B) {
				benchmarkParsing(b, parser, data.tokenString, jwt.MapClaims{})
			})
		case *jwt.StandardClaims:
			b.Run("standard_claims", func(b *testing.B) {
				benchmarkParsing(b, parser, data.tokenString, &jwt.StandardClaims{})
			})
		}
	}
}

// Helper method for benchmarking various parsing methods
func benchmarkParsing(b *testing.B, parser *jwt.Parser, tokenString string, claims jwt.Claims) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Helper method for benchmarking various signing methods
func benchmarkSigning(b *testing.B, method jwt.SigningMethod, key interface{}) {
	b.Helper()
	t := jwt.New(method)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := t.SignedString(key); err != nil {
				b.Fatal(err)
			}
		}
	})
}
