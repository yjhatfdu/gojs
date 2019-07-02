package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSContextRef.h>
// typedef bool (*JSShouldTerminateCallback) (JSContextRef ctx, void* context);
// void JSContextGroupSetExecutionTimeLimit(JSContextGroupRef group, double limit, void* callback, JSContextRef context);
// void JSContextGroupClearExecutionTimeLimit(JSContextGroupRef group);
// bool shouldTerminate(JSContextRef ctx, void* context){
// return 1;}
import "C"
import (
	"time"
	"unsafe"
)

// Context wraps a JavaScriptCore JSContextRef.
type Context struct {
	ref C.JSContextRef
}

// GlobalContext wraps a JavaScriptCore JSGlobalContextRef.
type GlobalContext Context

func NewContext() *Context {
	ctx := new(Context)

	c_nil := C.JSClassRef(unsafe.Pointer(uintptr(0)))
	ctx.ref = C.JSContextRef(C.JSGlobalContextCreate(c_nil))
	return ctx
}

type RawContext C.JSContextRef

type RawGlobalContext C.JSGlobalContextRef

func NewContextFrom(raw RawContext) *Context {
	ctx := new(Context)
	ctx.ref = C.JSContextRef(raw)
	return ctx
}

func NewGlobalContextFrom(raw RawGlobalContext) *GlobalContext {
	ctx := new(GlobalContext)
	ctx.ref = C.JSContextRef(raw)
	return ctx
}

func (ctx *Context) Retain() {
	C.JSGlobalContextRetain(ctx.ref)
}

func (ctx *Context) Release() {
	C.JSGlobalContextRelease(ctx.ref)
}

func (ctx *Context) GlobalObject() *Object {
	ret := C.JSContextGetGlobalObject(ctx.ref)
	return ctx.newObject(ret)
}

func (ctx *Context) SetTimeLimit(duration time.Duration) {
	group := C.JSContextGetGroup(ctx.ref)
	C.JSContextGroupSetExecutionTimeLimit(group, C.double(duration.Seconds()), nil, nil)
}

func (ctx *Context) ClearTimeLimit() {
	group := C.JSContextGetGroup(ctx.ref)
	C.JSContextGroupClearExecutionTimeLimit(group)
}
