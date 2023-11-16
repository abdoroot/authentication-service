package auth

import (
	"context"
	"testing"
)

func TestIsUserAuthorized(t *testing.T) {
	tc := []struct {
		token    string
		expected bool 
	}{
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAxNDc4OTQsInVzZXJfZW1haWwiOiJhYmQuMjAwOTMwQGdtYWlsLmNvbSIsInVzZXJfaWQiOiIxMCJ9.H3UK9fwf7WVnSbUrx_9Six6WTIjtQhUdYVukMOuAyDA", true}, //valid for three hours
		{"unvalidToken", false}, //unvalid
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAxMzk3MzEsInVzZXJfZW1haWwiOiJhYmQuMjAwOTMwQGdtYWlsLmNvbSIsInVzZXJfaWQiOiIxMCJ9.FhIGooPAnI7YMeWsoLlRoaXnfgKFrauxUsp1dbZz2HY", false}, //expired
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAxNDgwNjcsInVzZXJfZW1haWwiOiJhYmQuMjAwOTMwQGdtYWlsLmNvbSIsInVzZXJfaWQiOiIxMCJ9.pCRy2WoUpO7JhjI_HjyJyoEEyDA8uVPIGmAsLpAcTAE", true}, //valid for three hours
	}

	for _, v := range tc {
		ctx := context.WithValue(context.TODO(), "token", v.token)
		if _,got := IsUserAuthorizedWithClaim(ctx); got != v.expected {
			t.Errorf("got %v expected :%v ", got, v.expected)
		}
	}

}
