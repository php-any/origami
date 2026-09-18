package node

import "github.com/php-any/origami/data"

// callVM 从调用上下文取 VM。方法帧 / BoundContext 的内层 pooled Context
// 回收后 GetVM() 可能为 nil，此时回退到钉在 ClassValue 上的 VM 或 BoundThis。
func callVM(ctx data.Context) data.VM {
	if ctx == nil {
		return nil
	}
	if vm := ctx.GetVM(); vm != nil {
		return vm
	}
	switch c := ctx.(type) {
	case *data.ClassMethodContext:
		if c.ClassValue != nil {
			if vm := c.ClassValue.GetVM(); vm != nil {
				return vm
			}
		}
	case *data.ClassValue:
		return c.GetVM()
	case *data.BoundContext:
		if c.BoundThis != nil {
			if vm := c.BoundThis.GetVM(); vm != nil {
				return vm
			}
		}
		if c.Context != nil && c.Context != ctx {
			return callVM(c.Context)
		}
	}
	return nil
}
