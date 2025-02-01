package auth

import (
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {

	key:= headers.Get("Authorization")

	parts:= strings.SplitN(key, " ",2)
	if len(parts)!=2 || parts[0]!="ApiKey"{
		return "", nil
	}

    return parts[1], nil
}