package node

import "github.com/php-any/origami/data"

type TickStatement struct {
	Statement data.GetValue
	Interval  int
}

func (s *TickStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	value, ctl := s.Statement.GetValue(ctx)
	if ctl == nil {
		if provider, ok := ctx.GetVM().(data.TickProvider); ok {
			ctl = provider.PHPTickState().Advance(ctx, s.Interval)
		}
	}
	return value, ctl
}
