package auth

import "context"

func IsUserAuthorized(ctx context.Context) bool {
	token := ctx.Value("token")
	if _, ok := token.(string); ok {
		//do dome thin
	}

	return true
}
