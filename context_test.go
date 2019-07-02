package gojs

import (
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()
}

func TestContext2(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	ctx.Retain()
	defer ctx.Release()
}

func TestContextGlobalObject(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	obj := ctx.GlobalObject()
	if obj == nil {
		t.Errorf("ctx.GlobalObject() returned nil")
	}
	if obj.ToValue().Type() != TypeObject {
		t.Errorf("ctx.GlobalObject() did not return a javascript object")
	}
}

func TestNewContextWithTimeLimit(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()
	ctx.SetTimeLimit(time.Second)
	_, _ = ctx.EvaluateScript(`while(1){}`, nil, "", 0)

}
