package repo

import (
	"goboot/internal/repository/cond"
	"gorm.io/gorm"
)

type Saver interface {
	// Save 如果 gorm.Model 中 ID 不为0，则创建，否则更新实体所有字段，包括零值字段，参数 v 为实体对象指针
	Save(v any) (any, error)
}

type Creator interface {
	// CreateOne 新增记录，参数 v 必须是一个指针，返回带有 id 的实体的指针
	CreateOne(v any) (any, error)
	CreateOneTx(tx *gorm.DB, v any) (any, error)

	// CreateBatch 批量创建多条记录，参数 v 为实体 slice 或者实体 slice 的指针。创建完成后，slice 中的实体对象会自动填充主键 id。
	// 批量插入时会组装为一条 sql，如果出现错误，则返回 error
	CreateBatch(v any) error
	CreateBatchTx(tx *gorm.DB, v any) error

	// CreateBatchSize 批量创建多条记录，参数 v 为 slice 的指针。创建完成后，slice 中的实体对象会自动填充主键 id。
	// 该方法需要通过 size 参数指定每次批量创建的数量，如果 slice 的长度大于 size 或按照 size 大小切割为多个 sql
	CreateBatchSize(v any, size int) error
	CreateBatchSizeTx(tx *gorm.DB, v any, size int) error
}

type Updater interface {
	// UpdateField 按照 id 更新单个字段，参数 v 为实体指针，包含 id 主键
	UpdateField(v any, field string, val any) error
	UpdateFieldTx(tx *gorm.DB, v any, field string, val any) error

	// UpdateFieldCond 按照条件更新单个字段，参数 v 为实体对象指针
	UpdateFieldCond(v any, cond cond.Condition, field string, val any) error
	UpdateFieldCondTx(tx *gorm.DB, v any, cond cond.Condition, field string, val any) error

	// UpdateFields 根据实体对象更新对象中包含值的字段，参数 v 为实体对象指针，如果 cond 为 nil，
	// 则 v 必须包含 gorm.Model 的 ID 主键作为更新条件。
	// 如果 v 字段包含零值，则不更新，需要更新零值请使用 UpdateFieldsExact 方法
	UpdateFields(v any, cond cond.Condition) error
	UpdateFieldsTx(tx *gorm.DB, v any, cond cond.Condition) error
	UpdateFieldsSel(v any, cond cond.Condition, fields ...string) error
	UpdateFieldsSelTx(tx *gorm.DB, v any, cond cond.Condition, fields ...string) error

	// UpdateFieldsExact 根据实体对象更新 fields 指定的字段的值，包括零值，参数 v 为实体对象指针
	UpdateFieldsExact(v any, fields map[string]any, cond cond.Condition) error
	UpdateFieldsExactTx(tx *gorm.DB, v any, fields map[string]any, cond cond.Condition) error
}

type Deleter interface {
	// DeleteById 根据主键 id 删除记录，为逻辑删除，如果要采用物理删除，请使用 RemoveById 方法，参数 v 为实体对象指针
	DeleteById(v any, id uint) error
	DeleteByIdTx(tx *gorm.DB, v any, id uint) error

	// RemoveById 根据主键 id 删除记录，为物理删除，如果要采用逻辑删除，请使用 DeleteById 方法，参数 v 为实体对象指针
	RemoveById(v any, id uint) error
	RemoveByIdTx(tx *gorm.DB, v any, id uint) error
}

type Querier interface {
	// GetById 根据主键id查询实体，参数 v 为实体对象指针，并返回已填充值的实体数据指针，如果没有找到则返回 nil。
	// 该方法只查询未被逻辑删除的记录，如果需要查询已删除的记录，请使用 GetByIdExact 方法
	GetById(v any, id uint) (any, error)
	GetByIdTx(tx *gorm.DB, v any, id uint) (any, error)

	// GetByIdExact 根据主键 id 查询实体，参数 v 为实体对象指针，并返回已填充值的实体数据指针，如果没有找到则返回 nil。
	// 该方法可以查询已经被逻辑删除的数据，通常只需要使用 GetById 方法即可。
	GetByIdExact(v any, id uint) (any, error)
	GetByIdExactTx(tx *gorm.DB, v any, id uint) (any, error)

	// FindOne 根据传入的实体类型指针 v，查询符合条件 cond 的记录。确保符合条件的数据只有一条记录，如果需要查询多条记录，
	// 请使用 FindAll 方法。
	FindOne(v any, cond cond.Condition) (any, error)
	FindOneTx(tx *gorm.DB, v any, cond cond.Condition) (any, error)

	// Find 方法根据具体的实体类型指针 v, 查询符合条件 cond 的所有记录，返回 v 具体类型指针的切片。
	Find(v any, cond cond.Condition) (any, error)
	FindTx(tx *gorm.DB, v any, cond cond.Condition) (any, error)

	// FindCnt 方法根据具体的实体类型指针 v, 统计符合条件 cond 的所有记录的数量。
	FindCnt(v any, cond cond.Condition) (int64, error)
	FindCntTx(tx *gorm.DB, v any, cond cond.Condition) (int64, error)

	// FindRaw 使用原始 sql 查询。
	// 参数 v 指定返回的数据，有多条记录则为 slice，可以是指针也可以是 slice，如果确定只有一个记录则应该传递单个对象的结构指针，
	// 非指针结构该方法会自动转换。
	FindRaw(v any, sql string, args ...any) (any, error)
	FindRawTx(tx *gorm.DB, v any, sql string, args ...any) (any, error)
}

type Transactor interface {
	// Tx 在事务中执行函数 f
	Tx(f func(mgr Manager, tx *gorm.DB) error) error
}

type Manager interface {
	Saver
	Creator
	Updater
	Deleter
	Querier

	Transactor
}
