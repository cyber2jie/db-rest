package db

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
)

func TransformByScriptSource(row *DataRow, jsSource string) error {

	if row == nil {
		return errors.New("rowdata is nil")
	}
	_, fn, err := GetVm(jsSource)

	if err != nil {
		return err
	}
	fn(row)
	return nil
}

func GetVm(jsSource string) (*goja.Runtime, func(data *DataRow), error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	_, err := vm.RunString(jsSource)

	if err != nil {
		return nil, nil, err
	}

	_, ok := goja.AssertFunction(vm.Get(TRANSFORM_JS_FUNC_NAME))
	if !ok {
		return nil, nil, errors.New(fmt.Sprintf("%s is not a function", TRANSFORM_JS_FUNC_NAME))
	}

	var transformFn func(data *DataRow)

	vm.ExportTo(vm.Get(TRANSFORM_JS_FUNC_NAME), &transformFn)

	return vm, transformFn, nil
}
