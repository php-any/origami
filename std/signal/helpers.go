package signal

import (
	"errors"
	"os"
	"syscall"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

func parseSignal(v data.Value) (os.Signal, error) {
	if iv, ok := v.(data.AsInt); ok {
		n, err := iv.AsInt()
		if err != nil {
			return nil, err
		}
		return syscall.Signal(n), nil
	}
	return nil, errors.New("信号编号必须是 int")
}

func parseSignalsFromContext(ctx data.Context, startIndex int) ([]os.Signal, data.Control) {
	var sigs []os.Signal
	for i := startIndex; ; i++ {
		v, ok := ctx.GetIndexValue(i)
		if !ok || v == nil {
			break
		}
		if arr, ok := v.(*data.ArrayValue); ok {
			for _, item := range arr.ToValueList() {
				sig, err := parseSignal(item)
				if err != nil {
					return nil, utils.NewThrow(err)
				}
				sigs = append(sigs, sig)
			}
			continue
		}
		sig, err := parseSignal(v)
		if err != nil {
			return nil, utils.NewThrow(err)
		}
		sigs = append(sigs, sig)
	}
	return sigs, nil
}

func parseSignalsFromArray(v data.Value) ([]os.Signal, data.Control) {
	arr, ok := v.(*data.ArrayValue)
	if !ok {
		return nil, utils.NewThrow(errors.New("信号列表必须是 array"))
	}
	var sigs []os.Signal
	for _, item := range arr.ToValueList() {
		sig, err := parseSignal(item)
		if err != nil {
			return nil, utils.NewThrow(err)
		}
		sigs = append(sigs, sig)
	}
	return sigs, nil
}

func extractSignalChannel(v data.Value) (*SignalChannel, data.Control) {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil, utils.NewThrow(errors.New("参数必须是 Signal\\Channel 实例"))
	}
	sc, ok := cv.Class.(*SignalChannelClass)
	if !ok {
		return nil, utils.NewThrow(errors.New("参数必须是 Signal\\Channel 实例"))
	}
	if sc.channel == nil {
		return nil, utils.NewThrow(errors.New("Signal\\Channel 未初始化"))
	}
	return sc.channel, nil
}
