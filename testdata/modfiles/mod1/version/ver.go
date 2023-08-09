package version

import (
	"context"

	"github.com/LK4D4/joincontext"
)

func Ver() context.Context {
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	ctx2 := context.Background()

	joineCtx, _ := joincontext.Join(ctx1, ctx2)

	return joineCtx
}
