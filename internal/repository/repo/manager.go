package repo

import (
	"fmt"
	"github.com/gophero/goal/errorx"
	"goboot/internal/repository/cond"
	"reflect"

	"gorm.io/gorm"
)

var (
	RecordNotFoundErr = fmt.Errorf("record not found")
	InvalidRowErr     = fmt.Errorf("found unexpected rows")
)

// GormManager DB 通用操作
type GormManager struct {
	rdb *gorm.DB
	wdb *gorm.DB
}

// NewGormManager 创建 GormManager 实例，需要传入 gorm.DB 以便操作数据库
func NewGormManager(rdb *gorm.DB, wdb *gorm.DB) Manager {
	s := new(GormManager)
	s.rdb = rdb
	s.wdb = wdb
	return s
}

func (c *GormManager) clone(tx *gorm.DB) Manager {
	return NewGormManager(tx, tx)
}

func (c *GormManager) CreateOne(v any) (any, error) {
	return c.CreateOneTx(c.wdb, v)
}

func (c *GormManager) CreateOneTx(tx *gorm.DB, v any) (any, error) {
	r := tx.Create(v)
	if r.Error != nil {
		return v, r.Error
	}
	if r.RowsAffected == 0 {
		return v, errorx.Wrapf(InvalidRowErr, fmt.Sprintf("create one failed, row affected is 0 which expects = 1"))
	}
	return v, nil
}

func (c *GormManager) CreateBatch(v any) error {
	return c.CreateBatchTx(c.wdb, v)
}

func (c *GormManager) CreateBatchTx(tx *gorm.DB, v any) error {
	r := tx.Create(v)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected == 0 {
		return errorx.Wrapf(InvalidRowErr, fmt.Sprintf("create batch failed, row affected is 0 which expect > 0"))
	}
	return nil
}

func (c *GormManager) CreateBatchSize(v any, size int) error {
	return c.CreateBatchSizeTx(c.wdb, v, size)
}

func (c *GormManager) CreateBatchSizeTx(tx *gorm.DB, v any, size int) error {
	r := tx.CreateInBatches(v, size)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected == 0 {
		return errorx.Wrapf(InvalidRowErr, fmt.Sprintf("create batch failed, row affected is 0 which expect > 0"))
	}
	return nil
}

func (c *GormManager) Save(v any) (any, error) {
	return c.SaveTx(c.wdb, v)
}

func (c *GormManager) SaveTx(tx *gorm.DB, v any) (any, error) {
	r := tx.Save(v)
	if r.Error != nil {
		return v, r.Error
	}
	if r.RowsAffected == 0 {
		return v, errorx.Wrapf(InvalidRowErr, fmt.Sprintf("save failed, row affected is 0 which expect = 1"))
	}
	return v, nil
}

func (c *GormManager) UpdateField(v any, field string, val any) error {
	return c.UpdateFieldTx(c.wdb, v, field, val)
}

func (c *GormManager) UpdateFieldTx(tx *gorm.DB, v any, field string, val any) error {
	r := tx.Model(v).Update(field, val)
	if r.Error != nil {
		return tx.Error
	}
	return nil
}

func (c *GormManager) UpdateFieldCond(v any, cond cond.Condition, field string, val any) error {
	return c.UpdateFieldCondTx(c.wdb, v, cond, field, val)
}

func (c *GormManager) UpdateFieldCondTx(tx *gorm.DB, v any, cond cond.Condition, field string, val any) error {
	vo := reflect.ValueOf(v)
	// cond 为 nil 时 v 必须为指针，且包含 ID
	if vo.Kind() != reflect.Ptr && cond != nil {
		v = reflect.New(vo.Type()).Interface()
	}
	txv := tx.Model(v)
	if cond != nil {
		txv = txv.Where(cond.Exp(), cond.Args()...)
	}
	txv = txv.Update(field, val)
	if txv.Error != nil {
		return txv.Error
	}
	return nil
}

func (c *GormManager) UpdateFields(v any, cond cond.Condition) error {
	return c.UpdateFieldsTx(c.wdb, v, cond)
}

func (c *GormManager) UpdateFieldsTx(tx *gorm.DB, v any, cond cond.Condition) error {
	txv := tx.Model(v)
	// 条件为nil，则v必须含有主键id
	if cond != nil {
		txv = txv.Where(cond.Exp, cond.Args()...)
	}
	txv = txv.Updates(v)
	if txv.Error != nil {
		return txv.Error
	}
	return nil
}

func (c *GormManager) UpdateFieldsSel(v any, cond cond.Condition, fields ...string) error {
	return c.UpdateFieldsSelTx(c.wdb, v, cond, fields...)
}

func (c *GormManager) UpdateFieldsSelTx(tx *gorm.DB, v any, cond cond.Condition, fields ...string) error {
	txv := tx.Model(v)
	if len(fields) > 0 {
		txv = txv.Select(fields)
	}
	if cond != nil {
		txv = txv.Where(cond.Exp(), cond.Args()...)
	}
	txv = txv.Updates(v)
	if txv.Error != nil {
		return txv.Error
	}
	return nil
}

func (c *GormManager) UpdateFieldsExact(v any, fields map[string]any, cond cond.Condition) error {
	return c.UpdateFieldsExactTx(c.wdb, v, fields, cond)
}

func (c *GormManager) UpdateFieldsExactTx(tx *gorm.DB, v any, fields map[string]any, cond cond.Condition) error {
	txv := tx.Model(v)
	if cond != nil {
		txv = txv.Where(cond.Exp(), cond.Args()...)
	}
	txv = txv.Updates(fields)
	if txv.Error != nil {
		return txv.Error
	}
	return nil
}

func (c *GormManager) DeleteById(v any, id uint) error {
	return c.DeleteByIdTx(c.wdb, v, id)
}

func (c *GormManager) DeleteByIdTx(tx *gorm.DB, v any, id uint) error {
	r := tx.Where("id = ?", id).Delete(&v)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected == 0 {
		return errorx.Wrapf(InvalidRowErr, fmt.Sprintf("delete failed, row affected is 0 which expects = 1"))
	}
	return nil
}

func (c *GormManager) RemoveById(v any, id uint) error {
	return c.RemoveByIdTx(c.wdb, v, id)
}

func (c *GormManager) RemoveByIdTx(tx *gorm.DB, v any, id uint) error {
	// 执行原生sql进行删除
	table, _ := c.tableName(v)
	r := tx.Exec(fmt.Sprintf("delete from %s where id = ?", table), id)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected == 0 {
		return errorx.Wrapf(InvalidRowErr, fmt.Sprintf("remove failed, row affected is 0 which expects = 1"))
	}
	return nil
}

func (c *GormManager) GetById(v any, id uint) (any, error) {
	return c.GetByIdTx(c.rdb, v, id)
}

func (c *GormManager) GetByIdTx(tx *gorm.DB, v any, id uint) (any, error) {
	r := tx.Where("id = ?", id).Find(v)
	if r.Error != nil {
		return v, r.Error
	}
	return v, nil
}

func (c *GormManager) GetByIdExact(v any, id uint) (any, error) {
	return c.GetByIdExactTx(c.rdb, v, id)
}

func (c *GormManager) GetByIdExactTx(tx *gorm.DB, v any, id uint) (any, error) {
	// query table name
	table, _ := c.tableName(v)
	// execute raw sql
	r := tx.Raw(fmt.Sprintf("select * from %s where id = ?", table), id).Find(v)
	if r.Error != nil {
		return v, r.Error
	}
	return v, nil
}

// tableName 查询实体对应的表明，实体 v 必须是一个指针，返回表名称，如果出现错误，则返回 error
func (c *GormManager) tableName(v any) (string, error) {
	stmt := &gorm.Statement{DB: c.rdb}
	err := stmt.Parse(v)
	if err != nil {
		return "", err
	}
	table := stmt.Schema.Table
	return table, nil
}

func (c *GormManager) Tx(f func(s Manager, tx *gorm.DB) error) error {
	return c.wdb.Transaction(func(tx *gorm.DB) error {
		return f(c.clone(tx), tx)
	})
}

func (c *GormManager) FindOne(v any, cond cond.Condition) (any, error) {
	return c.FindOneTx(c.rdb, v, cond)
}

func (c *GormManager) FindOneTx(tx *gorm.DB, v any, cond cond.Condition) (any, error) {
	txv := tx.Model(v)
	if cond != nil {
		txv = txv.Where(cond.Exp(), cond.Args()...)
	}
	r := txv.Find(v)
	if r.Error != nil {
		return v, r.Error
	}
	return v, nil
}

func (c *GormManager) Find(v any, cond cond.Condition) (any, error) {
	return c.FindTx(c.rdb, v, cond)
}

func (c *GormManager) FindTx(tx *gorm.DB, v any, cond cond.Condition) (any, error) {
	txv := tx.Model(v)
	if cond != nil {
		txv = txv.Where(cond.Exp(), cond.Args()...)
	}
	tp := reflect.TypeOf(v) // 指针的具体类型
	rs := reflect.MakeSlice(reflect.SliceOf(tp), 0, 0).Interface()
	r := txv.Find(&rs)
	if r.Error != nil {
		return v, r.Error
	}
	return rs, nil
}

func (c *GormManager) FindCnt(v any, cond cond.Condition) (int64, error) {
	return c.FindCntTx(c.rdb, v, cond)
}

func (c *GormManager) FindCntTx(tx *gorm.DB, v any, cond cond.Condition) (int64, error) {
	txv := tx.Model(v)
	var cnt int64
	if cond != nil {
		txv = txv.Where(cond.Exp(), cond.Args()...)
	}
	r := txv.Count(&cnt)
	return cnt, r.Error
}

func (c *GormManager) FindRaw(v any, sql string, args ...any) (any, error) {
	return c.FindRawTx(c.rdb, v, sql, args...)
}

func (c *GormManager) FindRawTx(tx *gorm.DB, v any, sql string, args ...any) (any, error) {
	txv := tx.Raw(sql, args...)
	tp := reflect.ValueOf(v).Type()
	if tp.Kind() == reflect.Ptr {
		txv = txv.Scan(v)
	} else if tp.Kind() == reflect.Slice {
		txv = txv.Scan(&v)
	} else if tp.Kind() == reflect.Struct {
		// 如果传递结构体而非指针，则需要动态创建指针
		// 否则 panic: reflect: reflect.Value.Set using unaddressable value
		v = reflect.New(tp).Interface()
		txv = txv.Scan(v)
		// 返回解引用对象，非指针
		v = reflect.ValueOf(v).Elem().Interface()
	} else {
		txv = txv.Scan(v)
	}

	if txv.Error != nil {
		return v, txv.Error
	}
	return v, nil
}
