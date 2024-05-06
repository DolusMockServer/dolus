package schema

import (
	"github.com/labstack/echo/v4"
)

type Mapper interface {
	MapToRequestParameters(ctx echo.Context) RequestParameters
}

type MapperImpl struct{}

// check that MapperImpl implements Mapper
var _ Mapper = (*MapperImpl)(nil)

func NewMapper() *MapperImpl {
	return &MapperImpl{}
}

// Map implements Mapper.
func (*MapperImpl) MapToRequestParameters(ctx echo.Context) RequestParameters {
	return RequestParameters{
		PathParams:  getPathParams(ctx.ParamNames(), ctx.Param),
		QueryParams: ctx.QueryParams(),
	}
}

func getPathParams(
	paramNames []string,
	params func(string) string,
) (paramMap map[string]string) {
	if len(paramNames) > 0 {
		paramMap = make(map[string]string)
		for _, name := range paramNames {
			paramMap[name] = params(name)
		}
	}
	return
}
