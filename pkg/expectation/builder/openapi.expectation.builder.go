package builder

import (
	"strconv"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/generator"

	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/pkg/expectation/loader"
	"github.com/DolusMockServer/dolus/pkg/logger"
	"github.com/DolusMockServer/dolus/pkg/schema"
)

type OpenApiExpectationBuilder struct {
	loader         loader.Loader[loader.OpenAPISpecLoadType]
	fieldGenerator generator.Generator
}

var _ ExpectationBuilder = &OpenApiExpectationBuilder{}

func NewOpenApiExpectationBuilder(
	file string,
	fieldGenerator generator.Generator,
) *OpenApiExpectationBuilder {
	return &OpenApiExpectationBuilder{
		loader:         loader.NewOpenApiSpecLoader(file),
		fieldGenerator: fieldGenerator,
	}
}

func (oeb *OpenApiExpectationBuilder) BuildExpectations() ([]expectation.DolusExpectation, error) {
	if s, err := oeb.loader.Load(); err != nil {
		return nil, err
	} else {
		return oeb.buildExpectationsFromOpenApiSpec(s), nil
	}
}

func (oeb *OpenApiExpectationBuilder) buildExpectationsFromOpenApiSpec(
	spec *loader.OpenAPISpecLoadType,
) (expectations []expectation.DolusExpectation) {
	for path := range spec.Paths {
		for method, operation := range spec.Paths[path].Operations() {
			for code, ref := range operation.Responses {
				if path != "/store/order/{orderId}" || code != "200" {
					continue
				}
				// if p != "/" || code != "200" {
				// 	continue
				// }

				logger.Log.Info(path, " ", code)
				//  TODO check if uricionPath is needed here

				// requestSchema := schema.RequestSchemaFromOpenApi3RequestRef(ref)
				// engine must store for each path method code then check that

				status, _ := strconv.Atoi(code)

				expectations = append(expectations, expectation.DolusExpectation{
					Priority: 0,
					Request: expectation.DolusRequest{
						Body:        nil, // TODO: what should this be?
						OpenApiPath: path,
						Method:      method,
					},
					Response: expectation.DolusResponse{
						Body: dstruct.NewGeneratedStructWithConfig(
							schema.ResponseSchemaFromOpenApi3ResponseRef(
								ref,
								"application/json",
							),
							&oeb.fieldGenerator,
						),
						Status: status,
					},
				})

			}
		}
	}
	return
}
