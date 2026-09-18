package routing

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	nsException              = "Symfony\\Component\\Routing\\Exception\\"
	exceptionInterfaceName   = nsException + "ExceptionInterface"
	exInvalidArgument        = nsException + "InvalidArgumentException"
	exInvalidParameter       = nsException + "InvalidParameterException"
	exLogic                  = nsException + "LogicException"
	exMethodNotAllowed       = nsException + "MethodNotAllowedException"
	exMissingMandatory       = nsException + "MissingMandatoryParametersException"
	exNoConfiguration        = nsException + "NoConfigurationException"
	exResourceNotFound       = nsException + "ResourceNotFoundException"
	exRouteCircularReference = nsException + "RouteCircularReferenceException"
	exRouteNotFound          = nsException + "RouteNotFoundException"
	exRuntime                = nsException + "RuntimeException"
)

func newExceptionInterface() data.InterfaceStmt {
	return node.NewInterfaceStatement(nil, exceptionInterfaceName, []string{"Throwable"}, nil)
}

func newSimpleException(name, parent string, impl []string) *rtClass {
	c := newRtClass(name, strPtr(parent), impl, []data.Property{
		protProp("message", data.NewStringValue("")),
		protProp("code", data.NewIntValue(0)),
		protProp("file", data.NewStringValue("")),
		protProp("line", data.NewIntValue(0)),
		protProp("previous", data.NewNullValue()),
	})
	c.add(methDef("__construct",
		[]data.GetValue{
			param("message", 0, data.NewStringValue(""), nil),
			param("code", 1, data.NewIntValue(0), nil),
			param("previous", 2, data.NewNullValue(), nil),
		},
		[]data.Variable{
			variable("message", 0, nil),
			variable("code", 1, nil),
			variable("previous", 2, nil),
		},
		func(ctx data.Context) (data.GetValue, data.Control) {
			cv := rtSelf(ctx)
			setProp(cv, "message", data.NewStringValue(argString(ctx, 0, "")))
			setProp(cv, "code", data.NewIntValue(argInt(ctx, 1, 0)))
			if prev := arg(ctx, 2); prev != nil {
				setProp(cv, "previous", prev)
			}
			return data.NewNullValue(), nil
		}))
	c.add(meth("getMessage", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "message", "")), nil
	}))
	c.add(meth("getCode", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewIntValue(propInt(rtSelf(ctx), "code", 0)), nil
	}))
	c.add(meth("getFile", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "file", "")), nil
	}))
	c.add(meth("getLine", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewIntValue(propInt(rtSelf(ctx), "line", 0)), nil
	}))
	c.add(meth("getPrevious", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return prop(rtSelf(ctx), "previous"), nil
	}))
	c.add(meth("getTrace", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return emptyArray(), nil
	}))
	c.add(meth("getTraceAsString", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(""), nil
	}))
	c.add(meth("__toString", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		cv := rtSelf(ctx)
		return data.NewStringValue(cv.GetName() + ": " + propString(cv, "message", "")), nil
	}))
	return c
}

func newMethodNotAllowedException() *rtClass {
	c := newSimpleException(exMethodNotAllowed, "RuntimeException", []string{exceptionInterfaceName})
	c.props = append(c.props, protProp("allowedMethods", emptyArray()))
	c.methods["__construct"] = methDef("__construct",
		[]data.GetValue{
			param("allowedMethods", 0, emptyArray(), nil),
			param("message", 1, data.NewStringValue(""), nil),
			param("code", 2, data.NewIntValue(0), nil),
			param("previous", 3, data.NewNullValue(), nil),
		},
		[]data.Variable{
			variable("allowedMethods", 0, nil),
			variable("message", 1, nil),
			variable("code", 2, nil),
			variable("previous", 3, nil),
		},
		func(ctx data.Context) (data.GetValue, data.Control) {
			cv := rtSelf(ctx)
			methods := toStringSlice(arg(ctx, 0))
			for i := range methods {
				methods[i] = strings.ToUpper(methods[i])
			}
			setProp(cv, "allowedMethods", stringsToArray(methods))
			setProp(cv, "message", data.NewStringValue(argString(ctx, 1, "")))
			setProp(cv, "code", data.NewIntValue(argInt(ctx, 2, 0)))
			if prev := arg(ctx, 3); prev != nil {
				setProp(cv, "previous", prev)
			}
			return data.NewNullValue(), nil
		})
	c.add(meth("getAllowedMethods", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "allowedMethods"), nil
	}))
	return c
}

func newMissingMandatoryException() *rtClass {
	c := newSimpleException(exMissingMandatory, "InvalidArgumentException", []string{exceptionInterfaceName})
	c.props = append(c.props,
		privProp("routeName", data.NewStringValue("")),
		privProp("missingParameters", emptyArray()),
	)
	c.methods["__construct"] = methDef("__construct",
		[]data.GetValue{
			param("routeName", 0, data.NewStringValue(""), nil),
			param("missingParameters", 1, emptyArray(), nil),
			param("code", 2, data.NewIntValue(0), nil),
			param("previous", 3, data.NewNullValue(), nil),
		},
		[]data.Variable{
			variable("routeName", 0, nil),
			variable("missingParameters", 1, nil),
			variable("code", 2, nil),
			variable("previous", 3, nil),
		},
		func(ctx data.Context) (data.GetValue, data.Control) {
			cv := rtSelf(ctx)
			name := argString(ctx, 0, "")
			missing := toStringSlice(arg(ctx, 1))
			setProp(cv, "routeName", data.NewStringValue(name))
			setProp(cv, "missingParameters", stringsToArray(missing))
			msg := fmt.Sprintf("Some mandatory parameters are missing (\"%s\") to generate a URL for route \"%s\".",
				strings.Join(missing, "\", \""), name)
			setProp(cv, "message", data.NewStringValue(msg))
			setProp(cv, "code", data.NewIntValue(argInt(ctx, 2, 0)))
			if prev := arg(ctx, 3); prev != nil {
				setProp(cv, "previous", prev)
			}
			return data.NewNullValue(), nil
		})
	c.add(meth("getMissingParameters", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return propArray(rtSelf(ctx), "missingParameters"), nil
	}))
	c.add(meth("getRouteName", nil, func(ctx data.Context) (data.GetValue, data.Control) {
		return data.NewStringValue(propString(rtSelf(ctx), "routeName", "")), nil
	}))
	return c
}

func newRouteCircularException() *rtClass {
	c := newSimpleException(exRouteCircularReference, exRuntime, []string{exceptionInterfaceName})
	c.methods["__construct"] = methDef("__construct",
		[]data.GetValue{
			param("routeId", 0, nil, nil),
			param("path", 1, emptyArray(), nil),
		},
		[]data.Variable{
			variable("routeId", 0, nil),
			variable("path", 1, nil),
		},
		func(ctx data.Context) (data.GetValue, data.Control) {
			cv := rtSelf(ctx)
			id := argString(ctx, 0, "")
			path := strings.Join(toStringSlice(arg(ctx, 1)), " -> ")
			setProp(cv, "message", data.NewStringValue(
				fmt.Sprintf("Circular reference detected for route \"%s\", path: \"%s\".", id, path)))
			return data.NewNullValue(), nil
		})
	return c
}

func routingExceptionClasses() []data.ClassStmt {
	iface := []string{exceptionInterfaceName}
	return []data.ClassStmt{
		newSimpleException(exInvalidArgument, "InvalidArgumentException", iface),
		newSimpleException(exInvalidParameter, "InvalidArgumentException", iface),
		newSimpleException(exLogic, "LogicException", nil),
		newMethodNotAllowedException(),
		newMissingMandatoryException(),
		newSimpleException(exResourceNotFound, "RuntimeException", iface),
		newSimpleException(exNoConfiguration, exResourceNotFound, iface),
		newRouteCircularException(),
		newSimpleException(exRouteNotFound, "InvalidArgumentException", iface),
		newSimpleException(exRuntime, "RuntimeException", iface),
	}
}
