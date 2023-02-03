package dstruct

// func DoSchemasMatch(expectation, schema *DynamicStructModifierImpl) error {

// 	// for k, _ := range schema.allFields {
// 	// 	if expectation.allFields[k] == nil { // TODO need to check if field is required and that field type is the same
// 	// 		// TODO combine all errors into single field
// 	// 		return fmt.Errorf("schemas do not match field '%s' is missing from schema", k)
// 	// 	}
// 	// }

// 	for k, v := range expectation.allFields {
// 		if schema.allFields[k] == nil { // TODO need to check if field is required and that field type is the same
// 			//remove field
// 			expectation.DeleteField(k)
// 			fmt.Printf("%s %+v\n", k, v.structFieldName)
// 			return nil
// 		}
// 	}

// 	return nil
// }
