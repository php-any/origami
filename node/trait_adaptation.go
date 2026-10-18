package node

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
)

type traitMethodCandidate struct {
	Trait  string
	Method data.Method
}

// ComposeTraitMethods 按 PHP 规则把 trait 方法写入类或 trait。
// 类自己的方法优先，且能压过同名冲突；父类方法不能。
// insteadof 只从原方法名上排除列出的 trait，不删除其它名字上的别名。
// as 不改名时只调整该方法的可见性，不会换掉 insteadof 已经选定的实现。
func ComposeTraitMethods(vm data.VM, className string, instance, static map[string]data.Method, traitNames []string, aliases []data.TraitAlias, from data.From) data.Control {
	if instance == nil {
		instance = map[string]data.Method{}
	}
	if static == nil {
		static = map[string]data.Method{}
	}
	ownInstance, ownStatic := methodNameSet(instance), methodNameSet(static)
	if ctl := validateAliasTargets(vm, traitNames, aliases, from); ctl != nil {
		return ctl
	}
	instanceCandidates, ctl := collectTraitMethods(vm, traitNames, false)
	if ctl != nil {
		return ctl
	}
	staticCandidates, ctl := collectTraitMethods(vm, traitNames, true)
	if ctl != nil {
		return ctl
	}
	if ctl = excludePrecededMethods(vm, traitNames, aliases, instanceCandidates, false, from); ctl != nil {
		return ctl
	}
	if ctl = excludePrecededMethods(vm, traitNames, aliases, staticCandidates, true, from); ctl != nil {
		return ctl
	}
	if ctl = adaptSameNameMethods(vm, traitNames, aliases, instanceCandidates, false, from); ctl != nil {
		return ctl
	}
	if ctl = adaptSameNameMethods(vm, traitNames, aliases, staticCandidates, true, from); ctl != nil {
		return ctl
	}
	if ctl = insertTraitMethods(className, instance, instanceCandidates, ownInstance, from); ctl != nil {
		return ctl
	}
	if ctl = insertTraitMethods(className, static, staticCandidates, ownStatic, from); ctl != nil {
		return ctl
	}
	if ctl = insertTraitAliases(vm, className, instance, instanceCandidates, ownInstance, traitNames, aliases, false, from); ctl != nil {
		return ctl
	}
	return insertTraitAliases(vm, className, static, staticCandidates, ownStatic, traitNames, aliases, true, from)
}

func validateAliasTargets(vm data.VM, traitNames []string, aliases []data.TraitAlias, from data.From) data.Control {
	for _, alias := range aliases {
		if alias.Precedence {
			continue
		}
		if _, _, ok := lookupTraitMethod(vm, traitNames, alias, false); ok {
			continue
		}
		if _, _, ok := lookupTraitMethod(vm, traitNames, alias, true); ok {
			continue
		}
		return data.NewCompileFatal(from, fmt.Sprintf("An alias (%s) was defined for method %s(), but this method does not exist", alias.Alias, alias.Method))
	}
	return nil
}

func methodNameSet(methods map[string]data.Method) map[string]struct{} {
	names := make(map[string]struct{}, len(methods))
	for name := range methods {
		names[name] = struct{}{}
	}
	return names
}

func collectTraitMethods(vm data.VM, traitNames []string, static bool) (map[string][]traitMethodCandidate, data.Control) {
	candidates := map[string][]traitMethodCandidate{}
	for _, traitName := range traitNames {
		trait, ctl := vm.GetOrLoadClass(traitName)
		if ctl != nil {
			return nil, ctl
		}
		if trait == nil {
			return nil, data.NewCompileFatal(nil, fmt.Sprintf("Could not find trait %s", traitName))
		}
		if static {
			if cs, ok := trait.(*ClassStatement); ok {
				for name, method := range cs.StaticMethods {
					if method != nil {
						candidates[name] = append(candidates[name], traitMethodCandidate{Trait: traitName, Method: method})
					}
				}
			}
			continue
		}
		if cs, ok := trait.(*ClassStatement); ok {
			for name, method := range cs.Methods {
				if method != nil {
					candidates[name] = append(candidates[name], traitMethodCandidate{Trait: traitName, Method: method})
				}
			}
			continue
		}
		for _, method := range trait.GetMethods() {
			if method == nil {
				continue
			}
			candidates[method.GetName()] = append(candidates[method.GetName()], traitMethodCandidate{Trait: traitName, Method: method})
		}
	}
	return candidates, nil
}

func excludePrecededMethods(vm data.VM, traitNames []string, aliases []data.TraitAlias, candidates map[string][]traitMethodCandidate, static bool, from data.From) data.Control {
	for _, alias := range aliases {
		if !alias.Precedence {
			continue
		}
		winner, ok := matchingUsedTrait(traitNames, alias.Trait)
		if alias.Trait == "" || !ok {
			return data.NewCompileFatal(from, "Could not find trait "+alias.Trait)
		}
		if !traitProvidesMethod(vm, winner, alias.Method, static) {
			if traitProvidesMethod(vm, winner, alias.Method, !static) {
				continue
			}
			return data.NewCompileFatal(from, fmt.Sprintf("An alias (%s) was defined for method %s(), but this method does not exist", alias.Method, alias.Method))
		}
		excluded := make([]string, 0, len(alias.Excluded))
		for _, name := range alias.Excluded {
			resolved, ok := matchingUsedTrait(traitNames, name)
			if !ok {
				return data.NewCompileFatal(from, "Could not find trait "+name)
			}
			excluded = append(excluded, resolved)
		}
		list := candidates[alias.Method]
		kept := make([]traitMethodCandidate, 0, len(list))
		for _, candidate := range list {
			if traitListed(excluded, candidate.Trait) {
				continue
			}
			kept = append(kept, candidate)
		}
		if len(kept) == 0 {
			delete(candidates, alias.Method)
		} else {
			candidates[alias.Method] = kept
		}
	}
	return nil
}

func adaptSameNameMethods(vm data.VM, traitNames []string, aliases []data.TraitAlias, candidates map[string][]traitMethodCandidate, static bool, from data.From) data.Control {
	for _, alias := range aliases {
		if alias.Precedence || alias.Alias != alias.Method {
			continue
		}
		traitName, ok := matchingUsedTrait(traitNames, alias.Trait)
		if alias.Trait != "" && !ok {
			return data.NewCompileFatal(from, "Could not find trait "+alias.Trait)
		}
		if alias.Trait != "" && !traitProvidesMethod(vm, traitName, alias.Method, static) && !methodWasExcluded(aliases, traitName, alias.Method) {
			if !traitProvidesMethod(vm, traitName, alias.Method, !static) {
				return data.NewCompileFatal(from, fmt.Sprintf("An alias (%s) was defined for method %s(), but this method does not exist", alias.Alias, alias.Method))
			}
			continue
		}
		list := candidates[alias.Method]
		for i, candidate := range list {
			if alias.Trait != "" && candidate.Trait != traitName {
				continue
			}
			list[i].Method = AdaptTraitMethod(candidate.Method, alias)
		}
	}
	return nil
}

func methodWasExcluded(aliases []data.TraitAlias, traitName, method string) bool {
	for _, alias := range aliases {
		if !alias.Precedence || alias.Method != method {
			continue
		}
		if traitListed(alias.Excluded, traitName) {
			return true
		}
	}
	return false
}

func insertTraitMethods(className string, dest map[string]data.Method, candidates map[string][]traitMethodCandidate, own map[string]struct{}, from data.From) data.Control {
	for name, list := range candidates {
		if _, exists := own[name]; exists {
			continue
		}
		if len(list) > 1 {
			return data.NewCompileFatal(from, traitCollision(list[1].Trait, name, className, list[0].Trait))
		}
		if len(list) == 1 {
			dest[name] = list[0].Method
		}
	}
	return nil
}

func insertTraitAliases(vm data.VM, className string, dest map[string]data.Method, candidates map[string][]traitMethodCandidate, own map[string]struct{}, traitNames []string, aliases []data.TraitAlias, static bool, from data.From) data.Control {
	placed := map[string]string{}
	for _, alias := range aliases {
		if alias.Precedence || alias.Alias == alias.Method {
			continue
		}
		method, traitName, ok := lookupTraitMethod(vm, traitNames, alias, static)
		if !ok {
			if _, _, other := lookupTraitMethod(vm, traitNames, alias, !static); other {
				continue
			}
			return data.NewCompileFatal(from, fmt.Sprintf("An alias (%s) was defined for method %s(), but this method does not exist", alias.Alias, alias.Method))
		}
		if _, exists := own[alias.Alias]; exists {
			continue
		}
		if list := candidates[alias.Alias]; len(list) > 0 {
			return data.NewCompileFatal(from, traitCollision(traitName, alias.Alias, className, list[0].Trait))
		}
		if previous, exists := placed[alias.Alias]; exists {
			return data.NewCompileFatal(from, traitCollision(traitName, alias.Alias, className, previous))
		}
		dest[alias.Alias] = AdaptTraitMethod(method, alias)
		placed[alias.Alias] = traitName
	}
	return nil
}

func traitCollision(methodTrait, methodName, className, existingTrait string) string {
	return fmt.Sprintf("Trait method %s::%s has not been applied as %s::%s, because of collision with %s::%s", methodTrait, methodName, className, methodName, existingTrait, methodName)
}

func lookupTraitMethod(vm data.VM, traitNames []string, alias data.TraitAlias, static bool) (data.Method, string, bool) {
	candidates := traitNames
	if alias.Trait != "" {
		resolved, ok := matchingUsedTrait(traitNames, alias.Trait)
		if !ok {
			return nil, "", false
		}
		candidates = []string{resolved}
	}
	for _, traitName := range candidates {
		trait, ctl := vm.GetOrLoadClass(traitName)
		if ctl != nil || trait == nil {
			continue
		}
		if static {
			if cs, ok := trait.(*ClassStatement); ok {
				if method, ok := cs.StaticMethods[alias.Method]; ok && method != nil {
					return method, traitName, true
				}
			}
			if gsm, ok := trait.(data.GetStaticMethod); ok {
				if method, ok := gsm.GetStaticMethod(alias.Method); ok && method != nil {
					return method, traitName, true
				}
			}
			continue
		}
		if cs, ok := trait.(*ClassStatement); ok {
			if method, ok := cs.Methods[alias.Method]; ok && method != nil {
				return method, traitName, true
			}
		}
		for _, method := range trait.GetMethods() {
			if method != nil && method.GetName() == alias.Method {
				return method, traitName, true
			}
		}
	}
	return nil, "", false
}

func traitProvidesMethod(vm data.VM, traitName, method string, static bool) bool {
	_, _, ok := lookupTraitMethod(vm, []string{traitName}, data.TraitAlias{Trait: traitName, Method: method}, static)
	return ok
}

func matchingUsedTrait(traitNames []string, mentioned string) (string, bool) {
	if mentioned == "" {
		return "", false
	}
	for _, name := range traitNames {
		if name == mentioned || strings.HasSuffix(name, "\\"+mentioned) {
			return name, true
		}
	}
	return "", false
}

func traitListed(names []string, traitName string) bool {
	for _, name := range names {
		if name == traitName || strings.HasSuffix(traitName, "\\"+name) || strings.HasSuffix(name, "\\"+traitName) {
			return true
		}
	}
	return false
}

// AdaptTraitMethod copies declaration metadata; trait users must not mutate one another.
func AdaptTraitMethod(method data.Method, alias data.TraitAlias) data.Method {
	if original, ok := method.(*ClassMethod); ok {
		copy := *original
		copy.Name = alias.Alias
		if alias.Visibility != "" {
			copy.Modifier = data.NewModifier(alias.Visibility)
		}
		if alias.Final {
			copy.Flags |= data.MethodFinal
		}
		return &copy
	}
	return method
}
