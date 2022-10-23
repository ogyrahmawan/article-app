package middlewares

import (
	"errors"
	beego "github.com/beego/beego/v2/server/web"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"strings"
)

var (
	ErrMissingApiKey       = errors.New("missing api key")
	ErrInvalidApiKey       = errors.New("invalid api key")
	ErrApiKeyNotRegistered = errors.New("api key not registered")
	ErrApiKeyAuth          = errors.New("error api key authentication")
)

type (
	// KeyAuthConfig defines the config for KeyAuth middleware.
	KeyAuthConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper

		// KeyLookup is a string in the form of "<source>:<name>" that is used
		// to extract key from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "form:<name>"
		KeyLookup string `yaml:"key_lookup"`

		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Bearer".
		AuthScheme string

		// Validator is a function to validate key.
		// Required.
		Validator KeyAuthValidator

		// ErrorHandler defines a function which is executed for an invalid key.
		// It may be used to define a custom error.
		ErrorHandler KeyAuthErrorHandler
	}

	// KeyAuthValidator defines a function to validate KeyAuth credentials.
	KeyAuthValidator func(string, *beegoContext.Context) (bool, error)

	keyExtractor func(*beegoContext.Context) (string, error)

	// KeyAuthErrorHandler defines a function which is executed for an invalid key.
	KeyAuthErrorHandler func(error, *beegoContext.Context)
)

var (
	// DefaultKeyAuthConfig is the default KeyAuth middleware config.
	DefaultKeyAuthConfig = KeyAuthConfig{
		Skipper:    DefaultSkipper,
		KeyLookup:  "header:Authorization",
		AuthScheme: "Bearer",
	}
)

// KeyAuth returns an KeyAuth middleware.
//
// For valid key it calls the next handler.
// For invalid key, it sends "401 - Unauthorized" response.
// For missing key, it sends "400 - Bad Request" response.
func KeyAuth(fn KeyAuthValidator) beego.FilterChain {
	c := DefaultKeyAuthConfig
	c.Validator = fn
	return KeyAuthWithConfig(c)
}

// KeyAuthWithConfig returns an KeyAuth middleware with config.
// See `KeyAuth()`.
func KeyAuthWithConfig(config KeyAuthConfig) beego.FilterChain {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultKeyAuthConfig.Skipper
	}
	// Defaults
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultKeyAuthConfig.AuthScheme
	}
	if config.KeyLookup == "" {
		config.KeyLookup = DefaultKeyAuthConfig.KeyLookup
	}
	if config.Validator == nil {
		panic("key-auth middleware requires a validator function")
	}
	if config.ErrorHandler == nil {
		panic(" middleware requires a handler function")
	}
	// Initialize
	parts := strings.Split(config.KeyLookup, ":")
	extractor := keyFromHeader(parts[1], config.AuthScheme)
	switch parts[0] {
	case "query":
		extractor = keyFromQuery(parts[1])
	case "form":
		extractor = keyFromForm(parts[1])
	}

	return func(next beego.FilterFunc) beego.FilterFunc {
		return func(c *beegoContext.Context) {
			if config.Skipper(c) {
				next(c)
				return
			}

			// Extract and verify key
			key, err := extractor(c)
			if err != nil {
				config.ErrorHandler(err, c)
				return
			}
			// validate key
			valid, err := config.Validator(key, c)
			if err != nil {
				config.ErrorHandler(err, c)
				return
			} else if valid {
				next(c)
			} else {
				config.ErrorHandler(ErrApiKeyAuth, c)
				return
			}
		}
	}
}

// keyFromHeader returns a `keyExtractor` that extracts key from the request header.
func keyFromHeader(header string, authScheme string) keyExtractor {
	return func(c *beegoContext.Context) (string, error) {
		auth := c.Request.Header.Get(header)
		if auth == "" {
			return "", ErrMissingApiKey
		}
		if header == "Authorization" {
			l := len(authScheme)
			if len(auth) > l+1 && auth[:l] == authScheme {
				return auth[l+1:], nil
			}
			return "", ErrInvalidApiKey
		}
		return auth, nil
	}
}

// keyFromQuery returns a `keyExtractor` that extracts key from the query string.
func keyFromQuery(param string) keyExtractor {
	return func(c *beegoContext.Context) (string, error) {
		key := c.Input.Query(param)
		if key == "" {
			return "", ErrMissingApiKey
		}
		return key, nil
	}
}

// keyFromForm returns a `keyExtractor` that extracts key from the form.
func keyFromForm(param string) keyExtractor {
	return func(c *beegoContext.Context) (string, error) {
		key := c.Request.FormValue(param)
		if key == "" {
			return "", ErrMissingApiKey
		}
		return key, nil
	}
}
