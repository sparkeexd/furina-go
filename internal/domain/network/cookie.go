package network

import (
	"bytes"
	"fmt"
	"net/http"
)

const (
	// Token keys.
	ltokenV2Key = "ltoken_v2"
	ltmidV2Key  = "ltmid_v2"
	ltuidV2Key  = "ltuid_v2"
)

// Cookie data class that stores tokens required for calling HoYoLab endpoints.
type Cookie struct {
	ltokenV2 http.Cookie
	ltmidV2  http.Cookie
	ltuidV2  http.Cookie
}

// Constructor.
func NewCookie(ltokenV2 string, ltmidV2 string, ltuidV2 string) Cookie {
	return Cookie{
		ltokenV2: http.Cookie{Name: ltokenV2Key, Value: ltokenV2},
		ltmidV2:  http.Cookie{Name: ltmidV2Key, Value: ltmidV2},
		ltuidV2:  http.Cookie{Name: ltuidV2Key, Value: ltuidV2},
	}
}

// Converts cookie tokens into a string to be added into a request header.
// Each token is separated by semicolons.
// Similar to http.Request.AddCookie().
func (cookie Cookie) ParseCookie() string {
	count := 0
	tokens := cookie.Tokens()
	length := len(tokens)
	buffer := new(bytes.Buffer)

	for _, token := range tokens {
		count++
		fmt.Fprintf(buffer, "%s=%s", token.Name, token.Value)
		if count != length {
			fmt.Fprint(buffer, ";")
		}
	}

	return buffer.String()
}

// Return all cookie tokens as an array.
func (cookie Cookie) Tokens() []http.Cookie {
	return []http.Cookie{cookie.ltokenV2, cookie.ltmidV2, cookie.ltuidV2}
}
