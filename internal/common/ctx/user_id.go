package ctx

import "context"

type userIdCtx int

const UserIdCtxKey userIdCtx = 0

func GetUserId(ctx context.Context) int {
	return ctx.Value(UserIdCtxKey).(int)
}
