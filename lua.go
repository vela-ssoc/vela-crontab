package crontab

import (
	"errors"
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/audit"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"

	"reflect"
)

var (
	xEnv vela.Environment

	_CronTypeOf = reflect.TypeOf((*Cron)(nil)).String()
	invalidArgs = errors.New("invalid args , usage add(string , title , function)")
)

func (c *Cron) NewLuaTask(L *lua.LState) int {
	n := L.GetTop()
	if n != 3 {
		L.RaiseError("%v", invalidArgs)
		return 0
	}

	spec := L.CheckString(1)
	title := L.CheckString(2)
	pip := pipe.NewByLua(L, pipe.Env(xEnv), pipe.Seek(2))

	tab := lua.NewMap(2, false)
	tab.Set("spec", lua.S2L(spec))
	tab.Set("title", lua.S2L(title))
	code := L.CodeVM()

	eid, err := c.AddFunc(spec, func() {
		co := xEnv.Clone(L)
		//这里注意 多个函数同时触发
		pip.Do(tab, co, func(err error) {
			audit.Debug("%s %s 计划任务执行失败 ", title, spec).From(code).High().Put()
		})
		xEnv.Free(co)
	})

	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	c.masks = append(c.masks, newMask(spec, title))

	L.Push(lua.LNumber(eid))
	return 1
}

func (c *Cron) startL(L *lua.LState) int {
	xEnv.Start(L, c).From(c.CodeVM()).Do()
	return 0
}

func (c *Cron) Index(L *lua.LState, key string) lua.LValue {

	if key == "task" {
		return lua.NewFunction(c.NewLuaTask)
	}

	if key == "start" {
		return lua.NewFunction(c.startL)
	}

	return lua.LNil
}

func newLuaCron(L *lua.LState) int {
	name := auxlib.CheckProcName(L.Get(1), L)

	proc := L.NewVelaData(name, _CronTypeOf)
	if proc.IsNil() {
		cron := New(name)
		cron.CodeVM()
		proc.Set(New(name))
	} else {
		c := proc.Data.(*Cron)
		c.Close()
		c.name = name
		c.masks = c.masks[:0]
	}

	L.Push(proc)
	return 1
}

func WithEnv(env vela.Environment) {
	xEnv = env
	env.Set("crontab", lua.NewFunction(newLuaCron))
}
