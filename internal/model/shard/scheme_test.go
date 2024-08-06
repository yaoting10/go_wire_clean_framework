package shard

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type ITestModel interface {
	TestMethod()
}

type TestModel struct{}

func (t *TestModel) GetID() uint {
	panic("implement me")
}

func (t *TestModel) GetCreatedAt() *time.Time {
	panic("implement me")
}

func (t *TestModel) GetUpdatedAt() *time.Time {
	panic("implement me")
}

type TestModel1 struct {
	TestModel
}

func (t *TestModel1) TestMethod() {
}

func (t *TestModel1) GetID() uint {
	panic("implement me")
}

func (t *TestModel1) GetCreatedAt() *time.Time {
	panic("implement me")
}

func (t *TestModel1) GetUpdatedAt() *time.Time {
	panic("implement me")
}

type TestModel2 struct {
	TestModel
}

func (t *TestModel2) TestMethod() {
}

func (t *TestModel2) GetID() uint {
	panic("implement me")
}

func (t *TestModel2) GetCreatedAt() *time.Time {
	panic("implement me")
}

func (t *TestModel2) GetUpdatedAt() *time.Time {
	panic("implement me")
}

func TestName(t *testing.T) {
	SchemeManager.RegisterType("test.model", NewMetadata("test_model", &TestModel{}))
	SchemeManager.RegisterType("test.model1", NewMetadata("test_model1", &TestModel1{}))
	SchemeManager.RegisterType("test.model2", NewMetadata("test_model2", &TestModel2{}))
	if v, err := SchemeManager.New("test.model"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("%+v\n", reflect.TypeOf(v))
		_, ok := v.(ITestModel)
		fmt.Println(ok)
	}
	if v, err := SchemeManager.New("test.model1"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("%+v\n", reflect.TypeOf(v))
		_, ok := v.(TestModel)
		fmt.Println("testModel:", ok)
		_, ok = v.(ITestModel)
		fmt.Println("ITestModel:", ok)
	}
	if v, err := SchemeManager.New("test.model2"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("%+v\n", reflect.TypeOf(v))
		_, ok := v.(TestModel)
		fmt.Println("testModel:", ok)
		_, ok = v.(ITestModel)
		fmt.Println("ITestModel:", ok)
		fmt.Println(ok)
	}
}
