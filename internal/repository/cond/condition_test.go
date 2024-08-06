package cond_test

import (
	"fmt"
	"github.com/gophero/goal/testx"
	"goboot/internal/repository/cond"
	"testing"
)

func TestNewCond(t *testing.T) {
	c := cond.New("name like ? and age > ?", "张三%", 10)
	fmt.Println(c)
}

type User struct {
	Name    string
	Age     int
	Hobbies []string
}

func TestNewStructCond(t *testing.T) {
	c := cond.Parse(User{Name: "张三%", Age: 10})
	fmt.Println(c)
	c = cond.Parse(&User{Name: "张三%", Age: 10})
	fmt.Println(c)
	c = cond.Parse(&User{Name: "张三%", Age: 10, Hobbies: []string{"read", "write", "programming"}})
	fmt.Println(c)

	// qw := db.Wrapper{}
	// qw.Eq("name", "zhangsan")
}

type Case struct {
	cond *cond.Wrapper
	sql  string
}

func TestWrapCondIn(t *testing.T) {
	cs := []Case{
		{cond: cond.Wrap().In("name", "张三", "李四", "wangwu"), sql: "name in ('张三','李四','wangwu')"},
		{cond: cond.Wrap().In("age", 10, 20, 30, 40), sql: "age in (10,20,30,40)"},
	}

	l := testx.Wrap(t)
	l.Case("test condition")
	for _, c2 := range cs {
		s := c2.cond.String()
		want := c2.sql
		l.Require(s == want, "condition string (expect = actually): [%s] = [%s]", want, s)
	}
}

func TestCondWrapper(t *testing.T) {
	cases := []Case{
		{cond: cond.Wrap().Eq("name", "张三"), sql: "name = '张三'"},
		{cond: cond.Wrap().Eqf(func() bool { return true }, "name", "张三"), sql: "name = '张三'"},
		{cond: cond.Wrap().Eqf(func() bool { return false }, "name", "张三"), sql: ""},

		{cond: cond.Wrap().Ne("name", "张三"), sql: "name <> '张三'"},
		{cond: cond.Wrap().Nef(func() bool { return true }, "name", "张三"), sql: "name <> '张三'"},
		{cond: cond.Wrap().Nef(func() bool { return false }, "name", "张三"), sql: ""},

		{cond: cond.Wrap().Like("name", "张三"), sql: "name like '%张三%'"},
		{cond: cond.Wrap().LikeLeft("name", "张三"), sql: "name like '%张三'"},
		{cond: cond.Wrap().LikeRight("name", "张三"), sql: "name like '张三%'"},

		{cond: cond.Wrap().Lt("age", 100), sql: "age < 100"},
		{cond: cond.Wrap().Ltf(func() bool { return true }, "age", 100), sql: "age < 100"},
		{cond: cond.Wrap().Ltf(func() bool { return false }, "age", 100), sql: ""},

		{cond: cond.Wrap().Le("age", 100), sql: "age <= 100"},
		{cond: cond.Wrap().Lef(func() bool { return true }, "age", 100), sql: "age <= 100"},
		{cond: cond.Wrap().Lef(func() bool { return false }, "age", 100), sql: ""},

		{cond: cond.Wrap().Gt("age", 100), sql: "age > 100"},
		{cond: cond.Wrap().Gtf(func() bool { return true }, "age", 100), sql: "age > 100"},
		{cond: cond.Wrap().Gtf(func() bool { return false }, "age", 100), sql: ""},

		{cond: cond.Wrap().Ge("age", 100), sql: "age >= 100"},
		{cond: cond.Wrap().Gef(func() bool { return true }, "age", 100), sql: "age >= 100"},
		{cond: cond.Wrap().Gef(func() bool { return false }, "age", 100), sql: ""},

		{cond: cond.Wrap().Between("age", 18, 60), sql: "age between 18 and 60"},
		{cond: cond.Wrap().Betweenf(func() bool { return true }, "age", 18, 60), sql: "age between 18 and 60"},
		{cond: cond.Wrap().Betweenf(func() bool { return false }, "age", 18, 60), sql: ""},

		{cond: cond.Wrap().In("name"), sql: ""},
		{cond: cond.Wrap().In("name", "张三", "李四", "wangwu"), sql: "name in ('张三','李四','wangwu')"},
		{cond: cond.Wrap().In("age", 10, 20, 30, 40), sql: "age in (10,20,30,40)"},
	}

	l := testx.Wrap(t)
	l.Case("test condition")
	for _, c2 := range cases {
		s := c2.cond.String()
		want := c2.sql
		l.Require(s == want, "condition string (expect = actually): [%s] = [%s]", want, s)
	}
}
