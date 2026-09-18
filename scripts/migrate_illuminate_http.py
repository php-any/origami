#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import pathlib

def transform_request(text: str) -> str:
    text = text.replace("package httpfoundation", "package http", 1)
    text = text.replace('\t"net/http"\n', '\tnethttp "net/http"\n', 1)
    text = text.replace(
        '\t"github.com/php-any/origami/node"\n)',
        '\t"github.com/php-any/origami/node"\n\thttpfoundation "github.com/php-any/origami/std/symfony/http-foundation"\n)',
        1,
    )
    text = text.replace("*http.Request", "*nethttp.Request")
    text = text.replace("http.Request", "nethttp.Request")
    text = text.replace("http.CanonicalHeaderKey", "nethttp.CanonicalHeaderKey")
    text = text.replace("http.Method", "nethttp.Method")
    for name in [
        "NewInputBagValue",
        "NewParameterBagValue",
        "NewFileBagValue",
        "NewServerBagValue",
        "NewHeaderBagValue",
    ]:
        text = text.replace(name + "(", "httpfoundation." + name + "(")
    text = text.replace("headersFromRequest(", "httpfoundation.HeadersFromRequest(")
    return text


def transform_response(text: str) -> str:
    text = text.replace("package httpfoundation", "package http", 1)
    text = text.replace(
        '\t"github.com/php-any/origami/node"\n)',
        '\t"github.com/php-any/origami/node"\n\thttpfoundation "github.com/php-any/origami/std/symfony/http-foundation"\n)',
        1,
    )
    reps = [
        ("publicProp(", "httpfoundation.PublicProp("),
        ("pubMethod(", "httpfoundation.PubMethod("),
        ("param(", "httpfoundation.Param("),
        ("variable(", "httpfoundation.Variable("),
        ("fqnSymfonyResponse", "httpfoundation.FqnSymfonyResponse"),
        ("responseClassValue(", "httpfoundation.ResponseClassValue("),
        ("responseSelf(", "httpfoundation.ResponseSelf("),
        ("responseContent(", "httpfoundation.ResponseContent("),
        ("responseStatusCode(", "httpfoundation.ResponseStatusCode("),
        ("responseStatusText(", "httpfoundation.ResponseStatusText("),
        ("responseHeaders(", "httpfoundation.ResponseHeaders("),
        ("responseSetProp(", "httpfoundation.ResponseSetProp("),
        ("headersMapFromValue(", "httpfoundation.HeadersMapFromValue("),
        ("intParam(", "httpfoundation.IntParam("),
        ("boolParam(", "httpfoundation.BoolParam("),
        ("optionalStringParam(", "httpfoundation.OptionalStringParam("),
        ("createResponseHeaders(", "httpfoundation.CreateResponseHeaders("),
        ("applyStatusCode(", "httpfoundation.ApplyStatusCode("),
        ("respHeaderSet(", "httpfoundation.RespHeaderSet("),
        ("respHeaderRemove(", "httpfoundation.RespHeaderRemove("),
        ("throwNamed(", "httpfoundation.ThrowNamed("),
        ("callObjMethod(", "httpfoundation.CallObjMethod("),
        ("HeaderBagFrom(", "httpfoundation.HeaderBagFrom("),
        ("GetHeaderBagAll(", "httpfoundation.GetHeaderBagAll("),
        ("valueToAssocMap(", "httpfoundation.ValueToAssocMap("),
        ("ResponseHeaderBagFrom(", "httpfoundation.ResponseHeaderBagFrom("),
        ("setBagCookie(", "httpfoundation.SetBagCookie("),
        ("&BagCookie{", "&httpfoundation.BagCookie{"),
    ]
    for a, b in reps:
        text = text.replace(a, b)
    old = """func listMethodByName(list []data.Method, name string) (*bagMethod, bool) {
	for _, m := range list {
		if m.GetName() == name {
			if bm, ok := m.(*bagMethod); ok {
				return bm, true
			}
		}
	}
	return nil, false
}"""
    new = """func listMethodByName(list []data.Method, name string) (data.Method, bool) {
	for _, m := range list {
		if m.GetName() == name {
			return m, true
		}
	}
	return nil, false
}"""
    text = text.replace(old, new)
    text = text.replace(
        "m.modifier = data.ModifierProtected",
        "httpfoundation.SetMethodModifier(m, data.ModifierProtected)",
    )
    return text


def main() -> None:
    req = pathlib.Path("std/illuminate/http/request.go")
    resp = pathlib.Path("std/illuminate/http/response.go")
    req.write_text(transform_request(req.read_text(encoding="utf-8")), encoding="utf-8")
    resp.write_text(transform_response(resp.read_text(encoding="utf-8")), encoding="utf-8")
    print("ok", req.stat().st_size, resp.stat().st_size)


if __name__ == "__main__":
    main()
