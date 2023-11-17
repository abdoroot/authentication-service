package auth

import (
	"log"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestIsUserAuthorizedWithClaim(t *testing.T) {
	tc := []struct {
		id           int
		token        string
		expectedBool bool
		claims       jwt.MapClaims
	}{
		{1, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAxNDc4OTQsInVzZXJfZW1haWwiOiJhYmQuMjAwOTMwQGdtYWlsLmNvbSIsInVzZXJfaWQiOiIxMCJ9.H3UK9fwf7WVnSbUrx_9Six6WTIjtQhUdYVukMOuAyDA", true, nil},  //valid for three hours
		{2, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAxMzk3MzEsInVzZXJfZW1haWwiOiJhYmQuMjAwOTMwQGdtYWlsLmNvbSIsInVzZXJfaWQiOiIxMCJ9.FhIGooPAnI7YMeWsoLlRoaXnfgKFrauxUsp1dbZz2HY", false, nil}, //expired
		{3, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAxNDgwNjcsInVzZXJfZW1haWwiOiJhYmQuMjAwOTMwQGdtYWlsLmNvbSIsInVzZXJfaWQiOiIxMCJ9.pCRy2WoUpO7JhjI_HjyJyoEEyDA8uVPIGmAsLpAcTAE", true, nil},  //valid for three hours
		{4, "unvalidToken", false, nil}, //unvalid

	}

	for _, v := range tc {
		//ctx := context.WithValue(context.TODO(), "token", v.token)
		gotClaim, gotBool := IsUserAuthorizedWithClaim(v.token)
		if uId, ok := gotClaim["user_id"].(string); ok {
			log.Println(uId)
		}
		if gotBool != v.expectedBool {
			t.Errorf("%v: got %v expected :%v ", v.id, gotBool, v.expectedBool)
		}
	}

}

func TestRefreshToken(t *testing.T) {
}
