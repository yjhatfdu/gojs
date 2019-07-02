package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSObjectRef.h>
// #include "callback.h"
import "C"
import "unsafe"

type Object struct {
	ref C.JSObjectRef
	ctx *Context
}

func release_jsstringref_array(refs []C.JSStringRef) {
	for i := 0; i < len(refs); i++ {
		if refs[i] != nil {
			C.JSStringRelease(refs[i])
		}
	}
}

// Creates a new *Object given a C pointer to an JSObjectRef.
func (ctx *Context) newObject(ref C.JSObjectRef) *Object {
	obj := new(Object)
	obj.ref = ref
	obj.ctx = ctx
	return obj
}

func (ctx *Context) NewEmptyObject() *Object {
	obj := C.JSObjectMake(ctx.ref, nil, nil)
	return ctx.newObject(obj)
}

func (ctx *Context) NewObjectWithProperties(properties map[string]*Value) (*Object, error) {
	obj := ctx.NewEmptyObject()
	for name, val := range properties {
		err := obj.SetProperty(name, val, 0)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func (ctx *Context) NewArray(items []*Value) (*Object, error) {
	errVal := ctx.newErrorValue()

	ret := ctx.NewEmptyObject()
	if items != nil {
		carr, carrlen := ctx.newCValueArray(items)
		ret.ref = C.JSObjectMakeArray(ctx.ref, carrlen, carr, &errVal.ref)
	} else {
		ret.ref = C.JSObjectMakeArray(ctx.ref, 0, nil, &errVal.ref)
	}
	if errVal.ref != nil {
		return nil, errVal
	}
	return ret, nil
}

func (ctx *Context) NewDate() (*Object, error) {
	errVal := ctx.newErrorValue()

	ret := C.JSObjectMakeDate(ctx.ref,
		0, nil,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewDateWithMilliseconds(milliseconds float64) (*Object, error) {
	errVal := ctx.newErrorValue()

	param := ctx.NewNumberValue(milliseconds)

	ret := C.JSObjectMakeDate(ctx.ref,
		C.size_t(1), &param.ref,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewDateWithString(date string) (*Object, error) {
	errVal := ctx.newErrorValue()

	param := ctx.NewStringValue(date)

	ret := C.JSObjectMakeDate(ctx.ref,
		C.size_t(1), &param.ref,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewRegExp(regex string) (*Object, error) {
	errVal := ctx.newErrorValue()

	param := ctx.NewStringValue(regex)

	ret := C.JSObjectMakeRegExp(ctx.ref,
		C.size_t(1), &param.ref,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewRegExpFromValues(parameters []*Value) (*Object, error) {
	errVal := ctx.newErrorValue()

	ret := C.JSObjectMakeRegExp(ctx.ref,
		C.size_t(len(parameters)), &parameters[0].ref,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewFunction(name string, parameters []string, body string, source_url string, starting_line_number int) (*Object, error) {
	Cname := NewString(name)
	defer Cname.Release()
	paramLens := len(parameters)
	Cparameters := make([]C.JSStringRef, paramLens)
	defer release_jsstringref_array(Cparameters)
	for i := 0; i < paramLens; i++ {
		Cparameters[i] = (C.JSStringRef)(unsafe.Pointer(NewString(parameters[i])))
	}

	Cbody := NewString(body)
	defer Cbody.Release()

	var sourceRef *String
	if source_url != "" {
		sourceRef = NewString(source_url)
		defer sourceRef.Release()
	}
	var argsPointer *C.JSStringRef
	if paramLens > 0 {
		argsPointer = &Cparameters[0]
	}
	errVal := ctx.newErrorValue()
	ret := C.JSObjectMakeFunction(ctx.ref,
		(C.JSStringRef)(unsafe.Pointer(Cname)),
		C.unsigned(paramLens), argsPointer,
		(C.JSStringRef)(unsafe.Pointer(Cbody)),
		(C.JSStringRef)(unsafe.Pointer(sourceRef)),
		C.int(starting_line_number), &errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (obj *Object) GetPrototype() *Value {
	ret := C.JSObjectGetPrototype(obj.ctx.ref, obj.ref)
	return obj.ctx.newValue(ret)
}

func (obj *Object) SetPrototype(rhs *Value) {
	C.JSObjectSetPrototype(obj.ctx.ref, obj.ref, rhs.ref)
}

func (obj *Object) HasProperty(name string) bool {
	jsstr := NewString(name)
	defer jsstr.Release()

	ret := C.JSObjectHasProperty(obj.ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)))
	return bool(ret)
}

func (obj *Object) GetProperty(name string) (*Value, error) {
	jsstr := NewString(name)
	defer jsstr.Release()

	errVal := obj.ctx.newErrorValue()

	ret := C.JSObjectGetProperty(obj.ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)), &errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}

	return obj.ctx.newValue(ret), nil
}

func (obj *Object) GetPropertyAtIndex(index uint16) (*Value, error) {
	errVal := obj.ctx.newErrorValue()

	ret := C.JSObjectGetPropertyAtIndex(obj.ctx.ref, obj.ref, C.unsigned(index), &errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}

	return obj.ctx.newValue(ret), nil
}

func (obj *Object) SetProperty(name string, rhs *Value, attributes uint8) error {
	jsstr := NewString(name)
	defer jsstr.Release()

	errVal := obj.ctx.newErrorValue()

	C.JSObjectSetProperty(obj.ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)), rhs.ref,
		(C.JSPropertyAttributes)(attributes), &errVal.ref)
	if errVal.ref != nil {
		return errVal
	}

	return nil
}

func (obj *Object) SetPropertyAtIndex(index uint16, rhs *Value) error {
	errVal := obj.ctx.newErrorValue()

	C.JSObjectSetPropertyAtIndex(obj.ctx.ref, obj.ref, C.unsigned(index), rhs.ref, &errVal.ref)
	if errVal.ref != nil {
		return errVal
	}

	return nil
}

func (obj *Object) DeleteProperty(name string) (bool, error) {
	jsstr := NewString(name)
	defer jsstr.Release()

	errVal := obj.ctx.newErrorValue()

	ret := C.JSObjectDeleteProperty(obj.ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)), &errVal.ref)
	if errVal.ref != nil {
		return false, errVal
	}

	return bool(ret), nil
}

// ToValue returns the JSValueRef wrapper for the object.
//
// Any JSObjectRef can be safely cast to a JSValueRef.
// https://lists.webkit.org/pipermail/webkit-dev/2009-May/007530.html
func (obj *Object) ToValue() *Value {
	if obj == nil {
		panic("ToValue() called on nil *Object!")
	}
	return obj.ctx.newValue(C.JSValueRef(obj.ref))
}

func (obj *Object) IsFunction() bool {
	return bool(C.JSObjectIsFunction(obj.ctx.ref, obj.ref))
}

func (obj *Object) CallAsFunction(thisObject *Object, parameters []*Value) (*Value, error) {
	errVal := obj.ctx.newErrorValue()
	cParameters, n := obj.ctx.newCValueArray(parameters)
	if thisObject == nil {
		thisObject = obj.ctx.newObject(nil)
		//log.Println(thisObject.ref)
	}

	ret := C.JSObjectCallAsFunction(obj.ctx.ref, obj.ref, thisObject.ref, n, cParameters, &errVal.ref)

	if errVal.ref != nil {
		return nil, errVal
	}

	return obj.ctx.newValue(ret), nil
}

func (obj *Object) IsConstructor() bool {
	return bool(C.JSObjectIsConstructor(obj.ctx.ref, obj.ref))
}

func (obj *Object) CallAsConstructor(parameters []*Value) (*Value, error) {
	errVal := obj.ctx.newErrorValue()

	var Cparameters *C.JSValueRef
	if len(parameters) > 0 {
		// TODO: Is this safe?
		Cparameters = (*C.JSValueRef)(unsafe.Pointer(&parameters[0]))
	}

	ret := C.JSObjectCallAsConstructor(obj.ctx.ref, obj.ref,
		C.size_t(len(parameters)),
		Cparameters,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}

	return obj.ctx.newObject(ret).ToValue(), nil
}

//=========================================================
// PropertyNameArray
//

const (
	PropertyAttributeNone       = 0
	PropertyAttributeReadOnly   = 1 << 1
	PropertyAttributeDontEnum   = 1 << 2
	PropertyAttributeDontDelete = 1 << 3
)

const (
	ClassAttributeNone                 = 0
	ClassAttributeNoAutomaticPrototype = 1 << 1
)

type PropertyNameArray struct {
}

func (obj *Object) CopyPropertyNames() *PropertyNameArray {
	ret := C.JSObjectCopyPropertyNames(obj.ctx.ref, obj.ref)
	return (*PropertyNameArray)(unsafe.Pointer(ret))
}

func (ref *PropertyNameArray) Retain() {
	C.JSPropertyNameArrayRetain(C.JSPropertyNameArrayRef(unsafe.Pointer(ref)))
}

func (ref *PropertyNameArray) Release() {
	C.JSPropertyNameArrayRelease(C.JSPropertyNameArrayRef(unsafe.Pointer(ref)))
}

func (ref *PropertyNameArray) Count() uint16 {
	ret := C.JSPropertyNameArrayGetCount(C.JSPropertyNameArrayRef(unsafe.Pointer(ref)))
	return uint16(ret)
}

func (ref *PropertyNameArray) NameAtIndex(index uint16) string {
	jsstr := C.JSPropertyNameArrayGetNameAtIndex(C.JSPropertyNameArrayRef(unsafe.Pointer(ref)), C.size_t(index))
	defer C.JSStringRelease(jsstr)
	return (*String)(unsafe.Pointer(jsstr)).String()
}
