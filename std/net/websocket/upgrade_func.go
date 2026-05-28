package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type UpgradeFunction struct{}

func NewUpgradeFunction() data.FuncStmt { return &UpgradeFunction{} }

func (h *UpgradeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	request, err := utils.ConvertFromIndex[*http.Request](ctx, 0)
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	response, err := utils.ConvertFromIndex[http.ResponseWriter](ctx, 1)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	checkOrigin := false
	if v, ok := ctx.GetIndexValue(2); ok && v != nil {
		if bv, ok := v.(data.AsBool); ok {
			checkOrigin, _ = bv.AsBool()
		}
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if !checkOrigin {
				return true
			}
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}
			return origin == "http://"+r.Host || origin == "https://"+r.Host
		},
	}

	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	return data.NewProxyValue(NewConnClassFrom(conn), ctx.CreateBaseContext()), nil
}

func (h *UpgradeFunction) GetName() string            { return "Net\\Websocket\\upgrade" }
func (h *UpgradeFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *UpgradeFunction) GetIsStatic() bool          { return true }
func (h *UpgradeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "request", 0, nil, nil),
		node.NewParameter(nil, "response", 1, nil, nil),
		node.NewParameter(nil, "checkOrigin", 2, data.NewBoolValue(false), data.NewBaseType("bool")),
	}
}
func (h *UpgradeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "request", 0, nil),
		node.NewVariable(nil, "response", 1, nil),
		node.NewVariable(nil, "checkOrigin", 2, nil),
	}
}
func (h *UpgradeFunction) GetReturnType() data.Types { return data.NewBaseType("mixed") }
