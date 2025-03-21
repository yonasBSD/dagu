package digraph

import "strings"

// assertFunctions validates the function definitions.
func assertFunctions(fns []*funcDef) error {
	if fns == nil {
		return nil
	}

	nameMap := make(map[string]bool)
	for _, funcDef := range fns {
		if _, exists := nameMap[funcDef.Name]; exists {
			return wrapError("function", funcDef.Name, ErrDuplicateFunction)
		}
		nameMap[funcDef.Name] = true

		definedParamNames := strings.Split(funcDef.Params, " ")
		passedParamNames := extractParamNames(funcDef.Command)
		if len(definedParamNames) != len(passedParamNames) {
			return wrapError("function", funcDef.Name, ErrFuncParamsMismatch)
		}

		for i := 0; i < len(definedParamNames); i++ {
			if definedParamNames[i] != passedParamNames[i] {
				return wrapError("function", funcDef.Name, ErrFuncParamsMismatch)
			}
		}
	}

	return nil
}

// assertStepDef validates the step definition.
func assertStepDef(def stepDef, funcs []*funcDef) error {
	// Step name is required.
	if def.Name == "" {
		return wrapError("name", def.Name, ErrStepNameRequired)
	}

	// TODO: Validate executor config for each executor type.

	if def.Command == nil {
		if def.Executor == nil && def.Script == "" && def.Call == nil && def.Run == "" {
			return ErrStepCommandIsRequired
		}
	}

	// validate the function call if it exists.
	if def.Call != nil {
		calledFunc := def.Call.Function
		calledFuncDef := &funcDef{}
		for _, funcDef := range funcs {
			if funcDef.Name == calledFunc {
				calledFuncDef = funcDef
				break
			}
		}
		if calledFuncDef.Name == "" {
			return wrapError("function", calledFunc, ErrCallFunctionNotFound)
		}

		definedParamNames := strings.Split(calledFuncDef.Params, " ")
		if len(def.Call.Args) != len(definedParamNames) {
			return wrapError("function", calledFunc, ErrNumberOfParamsMismatch)
		}

		for _, paramName := range definedParamNames {
			_, exists := def.Call.Args[paramName]
			if !exists {
				return wrapError("function", calledFunc, ErrRequiredParameterNotFound)
			}
		}
	}

	return nil
}
