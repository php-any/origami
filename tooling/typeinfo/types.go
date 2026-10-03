// Package typeinfo describes analysis results without depending on PHP values
// or the execution VM. These types are never executable PHP declarations.
package typeinfo

import "strings"

type Type interface{ String() string }

type Inferred struct{ Types []Type }

func (*Inferred) String() string    { return "LspTypes" }
func (t *Inferred) Add(member Type) { t.Types = append(t.Types, member) }

type Tuple struct{ Types []Type }

func (t Tuple) String() string {
	parts := make([]string, len(t.Types))
	for i, member := range t.Types {
		parts[i] = member.String()
	}
	return strings.Join(parts, ", ")
}

// Generic retains type arguments for analysis. Runtime adapters may erase the
// arguments to the nominal class name, but PHP declarations cannot use them.
type Generic struct {
	Name  string
	Types []Type
}

func (t Generic) String() string { return t.Name }
