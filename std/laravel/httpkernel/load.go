package httpkernel

import "github.com/php-any/origami/data"

// Load leaves the framework Kernel to Composer. The native implementation does
// not yet have a complete class contract, so it cannot replace the official FQN.
// The HTTP host uses Resolve/Bootstrap/Sandbox/Handle/Terminate with either
// implementation and keeps process adaptation outside the PHP class.
func Load(_ data.VM) {}
