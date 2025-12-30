package node

type FuncYieldStackState struct {
	*FunctionStatement
	BodyIndex int
}
