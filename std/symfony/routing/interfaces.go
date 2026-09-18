package routing

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	requestContextAwareName = "Symfony\\Component\\Routing\\RequestContextAwareInterface"
	urlMatcherInterfaceName = "Symfony\\Component\\Routing\\Matcher\\UrlMatcherInterface"
	requestMatcherIfaceName = "Symfony\\Component\\Routing\\Matcher\\RequestMatcherInterface"
	redirectableMatcherName = "Symfony\\Component\\Routing\\Matcher\\RedirectableUrlMatcherInterface"
	urlGeneratorIfaceName   = "Symfony\\Component\\Routing\\Generator\\UrlGeneratorInterface"
	configurableReqName     = "Symfony\\Component\\Routing\\Generator\\ConfigurableRequirementsInterface"
	routeCompilerIfaceName  = "Symfony\\Component\\Routing\\RouteCompilerInterface"
	matcherDumperIfaceName  = "Symfony\\Component\\Routing\\Matcher\\Dumper\\MatcherDumperInterface"
	routerInterfaceName     = "Symfony\\Component\\Routing\\RouterInterface"

	refAbsoluteURL  = 0
	refAbsolutePath = 1
	refRelativePath = 2
	refNetworkPath  = 3
)

func routingInterfaces() []data.InterfaceStmt {
	ctxAware := node.NewInterfaceStatement(nil, requestContextAwareName, nil, []data.Method{
		node.NewInterfaceMethod(nil, "setContext", "public", []data.GetValue{param("context", 0, nil, nil)}, nil),
		node.NewInterfaceMethod(nil, "getContext", "public", nil, nil),
	})

	matcher := node.NewInterfaceStatement(nil, urlMatcherInterfaceName, []string{requestContextAwareName}, []data.Method{
		node.NewInterfaceMethod(nil, "match", "public", []data.GetValue{param("pathinfo", 0, nil, nil)}, data.NewBaseType("array")),
	})

	reqMatcher := node.NewInterfaceStatement(nil, requestMatcherIfaceName, nil, []data.Method{
		node.NewInterfaceMethod(nil, "matchRequest", "public", []data.GetValue{param("request", 0, nil, nil)}, data.NewBaseType("array")),
	})

	redir := node.NewInterfaceStatement(nil, redirectableMatcherName, nil, []data.Method{
		node.NewInterfaceMethod(nil, "redirect", "public", []data.GetValue{
			param("path", 0, nil, nil),
			param("route", 1, nil, nil),
			param("scheme", 2, data.NewNullValue(), nil),
		}, data.NewBaseType("array")),
	})

	gen := node.NewInterfaceStatement(nil, urlGeneratorIfaceName, []string{requestContextAwareName}, []data.Method{
		node.NewInterfaceMethod(nil, "generate", "public", []data.GetValue{
			param("name", 0, nil, nil),
			param("parameters", 1, emptyArray(), nil),
			param("referenceType", 2, data.NewIntValue(refAbsolutePath), nil),
		}, data.NewBaseType("string")),
	})
	gen.StaticProperty.Store("ABSOLUTE_URL", data.NewIntValue(refAbsoluteURL))
	gen.StaticProperty.Store("ABSOLUTE_PATH", data.NewIntValue(refAbsolutePath))
	gen.StaticProperty.Store("RELATIVE_PATH", data.NewIntValue(refRelativePath))
	gen.StaticProperty.Store("NETWORK_PATH", data.NewIntValue(refNetworkPath))

	cfg := node.NewInterfaceStatement(nil, configurableReqName, nil, []data.Method{
		node.NewInterfaceMethod(nil, "setStrictRequirements", "public", []data.GetValue{param("enabled", 0, nil, nil)}, nil),
		node.NewInterfaceMethod(nil, "isStrictRequirements", "public", nil, nil),
	})

	compiler := node.NewInterfaceStatement(nil, routeCompilerIfaceName, nil, []data.Method{
		node.NewInterfaceMethod(nil, "compile", "public", []data.GetValue{param("route", 0, nil, nil)}, nil),
	})

	dumper := node.NewInterfaceStatement(nil, matcherDumperIfaceName, nil, []data.Method{
		node.NewInterfaceMethod(nil, "dump", "public", []data.GetValue{param("options", 0, emptyArray(), nil)}, data.NewBaseType("string")),
		node.NewInterfaceMethod(nil, "getRoutes", "public", nil, nil),
	})

	router := node.NewInterfaceStatement(nil, routerInterfaceName, []string{urlMatcherInterfaceName, urlGeneratorIfaceName}, []data.Method{
		node.NewInterfaceMethod(nil, "getRouteCollection", "public", nil, nil),
	})

	return []data.InterfaceStmt{ctxAware, matcher, reqMatcher, redir, gen, cfg, compiler, dumper, router}
}
