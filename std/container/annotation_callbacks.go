package container

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

func bindClassAnnotation(ctx data.Context) data.Control {
	abstract, acl := annotationStringArg(ctx, 0)
	if acl != nil {
		return acl
	}
	if abstract == "" {
		return utils.NewThrow(errors.New("Bind 缺少 abstract 参数"))
	}
	cls, acl := annotationTargetClass(ctx)
	if acl != nil {
		return acl
	}
	if e := activeEngine(ctx); e != nil {
		e.Bind(abstract, cls.Name)
	}
	return nil
}

func injectParameterAnnotation(ctx data.Context) data.Control {
	service, _ := annotationStringArg(ctx, 0)
	param, className, acl := annotationTargetParameter(ctx)
	if acl != nil {
		return acl
	}
	if className == "" {
		return utils.NewThrow(errors.New("Container\\Annotation\\Inject 缺少所属类信息"))
	}
	metadataMarkConstructorInject(className, param.Index, param.Name, service, true)
	return nil
}

func namedParameterAnnotation(ctx data.Context) data.Control {
	name, acl := annotationStringArg(ctx, 0)
	if acl != nil {
		return acl
	}
	if name == "" {
		return utils.NewThrow(errors.New("Container\\Annotation\\Named 缺少 name 参数"))
	}
	param, className, acl := annotationTargetParameter(ctx)
	if acl != nil {
		return acl
	}
	if className == "" {
		return utils.NewThrow(errors.New("Container\\Annotation\\Named 缺少所属类信息"))
	}
	metadataMarkConstructorInject(className, param.Index, param.Name, name, false)
	return nil
}
