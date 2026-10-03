package httpfoundation

import "github.com/php-any/origami/data"

func init() {
	data.RegisterNativeRequestPolicy[*ParamBagData](data.NativeClone, func(scope *data.RequestObjectScope, source *ParamBagData) *ParamBagData {
		clone := source.clone()
		scope.RememberNativeState(source, clone)
		for key, value := range clone.values {
			clone.values[key] = scope.Bind(value)
		}
		return clone
	})
	data.RegisterNativeRequestPolicy[*HeaderBagData](data.NativeClone, func(scope *data.RequestObjectScope, source *HeaderBagData) *HeaderBagData { return source.clone() })
	data.RegisterNativeRequestPolicy[*ResponseHeaderBagData](data.NativeClone, func(scope *data.RequestObjectScope, source *ResponseHeaderBagData) *ResponseHeaderBagData {
		if source == nil {
			return nil
		}
		clone := newResponseHeaderBagData()
		clone.HeaderBagData = *source.HeaderBagData.clone()
		source.mu.RLock()
		defer source.mu.RUnlock()
		for key, name := range source.headerNames {
			clone.headerNames[key] = name
		}
		for key, value := range source.computedCacheControl {
			clone.computedCacheControl[key] = value
		}
		for domain, paths := range source.cookies {
			clone.cookies[domain] = make(map[string]map[string]*BagCookie, len(paths))
			for path, cookies := range paths {
				clone.cookies[domain][path] = make(map[string]*BagCookie, len(cookies))
				for name, cookie := range cookies {
					copy := *cookie
					copy.Value, copy.Domain, copy.SameSite = copyStringPointer(cookie.Value), copyStringPointer(cookie.Domain), copyStringPointer(cookie.SameSite)
					clone.cookies[domain][path][name] = &copy
				}
			}
		}
		return clone
	})
}

func copyStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}
