package hooks

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

func HookAPi() {
	fmt.Println(errors.UnimplementedErrorHint)
}
