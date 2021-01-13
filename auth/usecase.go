package auth

import "context"

const CtxUserKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, username string) error
	SignIn(ctx context.Context, username string) error
}
