package cond

import (
	"fmt"
	"github.com/gophero/goal/logx"
	"github.com/gophero/goal/stringx"
	"reflect"
	"strconv"
	"strings"
)

type Condition interface {
	Exp() string
	Args() []any
}

// condImpl is a struct of condition implementation, with a string expression and multi arguments.
type condImpl struct {
	exp  string
	args []any
}

func (d *condImpl) Exp() string {
	return d.exp
}

func (d *condImpl) Args() []any {
	return d.args
}

func (d *condImpl) String() string {
	return fmt.Sprintf("exp: %s, args: %v", d.exp, d.args)
}

func New(exp string, arg ...any) Condition {
	cond := new(condImpl)
	cond.exp = exp
	cond.args = append(cond.args, arg...)
	return cond
}

func Parse(v any) Condition {
	val := reflect.Indirect(reflect.ValueOf(v))
	typ := reflect.TypeOf(v)
	// 指针指向的类型
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	vs := make([]any, 0)
	cs := strings.Builder{}
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)       // value 对应的 field，可以获取值
		fn := typ.Field(i).Name // type 对应的 field，可以获取名称
		fv := f.Interface()     // field 的值
		switch f.Kind() {
		case reflect.Bool:
			if fv == true {
				cs.WriteString(equalsSql(fn))
				vs = append(vs, fv)
			}
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			if fv != 0 {
				cs.WriteString(equalsSql(fn))
				vs = append(vs, fv)
			}
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			if v != nil {
				v := reflect.ValueOf(fv)
				if v.Len() > 0 {
					cs.WriteString(inSql(fn))
					vs = append(vs, fv)
				}
			}
		case reflect.String:
			if fv != "" {
				cs.WriteString(equalsSql(fn))
				vs = append(vs, fv)
			}
		default:
			logx.Default.Warn("unsupported field type: %v.", f.Kind())
			continue
		}
	}
	if cs.String() == "" || len(vs) == 0 {
		panic("no field with value are found, but it is expected.")
	}
	cond := new(condImpl)
	cond.exp = strings.TrimSuffix(cs.String(), stringx.Space+stringx.And+stringx.Space)
	cond.args = vs
	logx.Default.Debug("parsed condition exp: %s, args: %v", cond.exp, cond.args)
	return cond
}

func inSql(fn string) string {
	sb := strings.Builder{}
	sb.WriteString(stringx.CamelCaseToUnderscore(fn))
	sb.WriteString(stringx.Space)
	sb.WriteString(stringx.In)
	sb.WriteString(stringx.Space)
	sb.WriteString(stringx.LeftBracket)
	sb.WriteString(stringx.QuestionMark)
	sb.WriteString(stringx.RightBracket)
	sb.WriteString(stringx.Space)
	sb.WriteString(stringx.And)
	sb.WriteString(stringx.Space)
	return sb.String()
}

func equalsSql(fn string) string {
	sb := strings.Builder{}
	// TODO 可以选择是否将驼峰名转为下划线
	sb.WriteString(stringx.CamelCaseToUnderscore(fn))
	sb.WriteString(stringx.Space)
	sb.WriteString(stringx.Equals)
	sb.WriteString(stringx.Space)
	sb.WriteString(stringx.QuestionMark)
	sb.WriteString(stringx.Space)
	sb.WriteString(stringx.And)
	sb.WriteString(stringx.Space)
	return sb.String()
}

// ==============================
// condition wrapper
// ==============================

func Wrap() *Wrapper {
	return new(Wrapper)
}

type Filter func() bool

// Wrapper is a wrapper of Condition which provides many convenient methods to combine sql query condition.
type Wrapper struct {
	expBuilder strings.Builder
	args       []any
}

func (cw *Wrapper) Exp() string {
	return strings.TrimSuffix(cw.expBuilder.String(), " and ")
}

func (cw *Wrapper) Args() []any {
	return cw.args
}

func (cw *Wrapper) Reset() *Wrapper {
	cw.expBuilder.Reset()
	cw.args = make([]any, 0)
	return cw
}

func (cw *Wrapper) String() string {
	exp := cw.Exp()
	for _, c := range cw.args {
		v := reflect.ValueOf(c)
		switch v.Kind() {
		case reflect.Bool:
			exp = strings.Replace(exp, "?", strconv.FormatBool(c.(bool)), 1)
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			exp = strings.Replace(exp, "?", strconv.Itoa(c.(int)), 1)
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			s := printSlice(c)
			exp = strings.Replace(exp, "?", s, 1)
		case reflect.String:
			exp = strings.Replace(exp, "?", "'"+c.(string)+"'", 1)
		default:
			panic(fmt.Sprintf("unsupported type: %T", c))
		}
	}
	return exp
}

func printSlice(c any) string {
	vs, ok := c.([]interface{})
	if !ok {
		panic(fmt.Sprintf("unsupported type: %T", c))
	}
	var sb strings.Builder
	for _, v := range vs {
		vv := reflect.ValueOf(v)
		switch vv.Kind() {
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			sb.WriteString(strconv.Itoa(v.(int)))
			sb.WriteRune(',')
		case reflect.String:
			sb.WriteString("'" + v.(string) + "'")
			sb.WriteRune(',')
		default:
			panic(fmt.Sprintf("unsupported element type of array or slice: %T.", v))
		}
	}
	return strings.TrimSuffix(sb.String(), ",")
}

func (cw *Wrapper) appendTo(filter Filter, exp string, v ...any) *Wrapper {
	if filter == nil || filter() {
		cw.expBuilder.WriteString(exp)
		cw.args = append(cw.args, v...)
		cw.and()
	}
	return cw
}

func (cw *Wrapper) and() *Wrapper {
	cw.expBuilder.WriteString(" and ")
	return cw
}

func (cw *Wrapper) or() *Wrapper {
	cw.expBuilder.WriteString(" or ")
	return cw
}

func (cw *Wrapper) Eq(column string, v any) *Wrapper {
	return cw.appendTo(nil, column+" = ?", v)
}

func (cw *Wrapper) Eqf(filter Filter, column string, v any) *Wrapper {
	return cw.appendTo(filter, column+" = ?", v)
}

func (cw *Wrapper) Ne(column string, v any) *Wrapper {
	return cw.appendTo(nil, column+" <> ?", v)
}

func (cw *Wrapper) Nef(filter Filter, column string, v any) *Wrapper {
	return cw.appendTo(filter, column+" <> ?", v)
}

func (cw *Wrapper) Gt(column string, v any) *Wrapper {
	return cw.appendTo(nil, column+" > ?", v)
}

func (cw *Wrapper) Gtf(filter Filter, column string, v any) *Wrapper {
	return cw.appendTo(filter, column+" > ?", v)
}

func (cw *Wrapper) Ge(column string, v any) *Wrapper {
	return cw.appendTo(nil, column+" >= ?", v)
}

func (cw *Wrapper) Gef(filter Filter, column string, v any) *Wrapper {
	return cw.appendTo(filter, column+" >= ?", v)
}

func (cw *Wrapper) Lt(column string, v any) *Wrapper {
	return cw.appendTo(nil, column+" < ?", v)
}

func (cw *Wrapper) Ltf(filter Filter, column string, v any) *Wrapper {
	return cw.appendTo(filter, column+" < ?", v)
}

func (cw *Wrapper) Le(column string, v any) *Wrapper {
	return cw.appendTo(nil, column+" <= ?", v)
}

func (cw *Wrapper) Lef(filter Filter, column string, v any) *Wrapper {
	return cw.appendTo(filter, column+" <= ?", v)
}

func (cw *Wrapper) Like(column string, v string) *Wrapper {
	return cw.appendTo(func() bool { return v != "" }, column+" like ?", "%"+v+"%")
}

func (cw *Wrapper) LikeLeft(column string, v string) *Wrapper {
	return cw.appendTo(func() bool { return v != "" }, column+" like ?", "%"+v)
}

func (cw *Wrapper) LikeRight(column string, v string) *Wrapper {
	return cw.appendTo(func() bool { return v != "" }, column+" like ?", v+"%")
}

func (cw *Wrapper) Between(column string, v1 any, v2 any) *Wrapper {
	return cw.appendTo(nil, column+" between ? and ?", v1, v2)
}

func (cw *Wrapper) Betweenf(filter Filter, column string, v1 any, v2 any) *Wrapper {
	return cw.appendTo(filter, column+" between ? and ?", v1, v2)
}

func (cw *Wrapper) In(column string, v ...any) *Wrapper {
	return cw.appendTo(func() bool { return v != nil && len(v) > 0 }, column+" in (?)", v)
}
