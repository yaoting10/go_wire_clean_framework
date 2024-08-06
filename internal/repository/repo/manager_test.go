package repo_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"goboot/internal/repository/cond"
	"goboot/internal/repository/repo"
	"goboot/test"
	"gorm.io/gorm"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	NilErr = fmt.Errorf("runtime error: invalid memory address or nil pointer dereference")

	// model需要自己创建并设置，否则时间会不匹配
	u1 = &user{Name: "zhangsan", Model: gorm.Model{CreatedAt: time.Now(), UpdatedAt: time.Now(), ID: 1}}
	u2 = &user{Name: "lisi", Model: gorm.Model{CreatedAt: time.Now(), UpdatedAt: time.Now(), ID: 2}}
)

type user struct {
	gorm.Model
	Name string `gorm:"varchar(64);default:''"`
}

func setupRepository() (repo.Manager, *test.MockGorm) {
	db := test.NewSqlmock().Gorm()
	mgr := repo.NewGormManager(db.DB, db.DB)
	return mgr, db
}

func TestMain(m *testing.M) {
	fmt.Println("begin to test manager")
	code := m.Run()
	fmt.Println("test manager end")
	os.Exit(code)
}

func TestGormManager_CreateBatch(t *testing.T) {
	mgr, mock := setupRepository()
	defer assert.NoError(t, mock.ExpectationsWereMet()) // 最后调用此方法，所有的期望值都匹配

	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		before  func()
		wantErr bool
	}{
		// Add test cases.
		// 传递空的slice会出错，所以wantErr为true
		{name: "no_data", args: args{[]*user{}}, wantErr: true, before: func() {
			mock.ExpectBegin()    // 数据库调用mock开始
			mock.ExpectRollback() // 设置期望值，期望回滚。空slice执行会失败，数据库会回滚
		}},
		{name: "create_1", args: args{[]*user{u1}}, wantErr: false, before: func() {
			mock.ExpectBegin()
			// 期望设置参数并成功执行sql，最后事务会提交
			mock.ExpectExec("INSERT INTO `users`").
				WithArgs(u1.CreatedAt, u1.UpdatedAt, u1.DeletedAt, u1.Name, u1.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
		}},
		{name: "create_2", args: args{[]*user{u1, u2}}, wantErr: false, before: func() {
			mock.ExpectBegin()
			// 期望设置参数并成功执行sql，最后事务会提交
			mock.ExpectExec("INSERT INTO `users`").
				// TODO 无法设置参数，自增主键ID批量时不会自动设置，可能是bug
				//WithArgs(u1.CreatedAt, u1.UpdatedAt, u1.DeletedAt, u1.Name, u2.CreatedAt, u2.UpdatedAt, u2.DeletedAt, u2.Name).
				WillReturnResult(sqlmock.NewResult(2, 1))
			mock.ExpectCommit()
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			if err := mgr.CreateBatch(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("CreateBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_CreateBatchSize(t *testing.T) {
	type args struct {
		v    any
		size int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		before  func() repo.Manager
	}{
		// Add test cases.
		// 被测试方法支持传递空slice而不会返回error，故wantErr是false
		{name: "no_data", args: args{v: []*user{}, size: 0}, wantErr: false, before: func() repo.Manager {
			mgr, mock := setupRepository()
			defer assert.NoError(t, mock.ExpectationsWereMet())
			mock.ExpectBegin()    // 数据库调用mock开始
			mock.ExpectRollback() // 设置期望值，期望回滚。空slice执行会失败，数据库会回滚
			return mgr
		}},
		{name: "create_1", args: args{v: []*user{u1}, size: 1}, wantErr: false, before: func() repo.Manager {
			mgr, mock := setupRepository()
			defer assert.NoError(t, mock.ExpectationsWereMet())
			mock.ExpectBegin()
			// 期望设置参数并成功执行sql，最后事务会提交
			mock.ExpectExec("INSERT INTO `users`").WithArgs(u1.CreatedAt, u1.UpdatedAt, u1.DeletedAt, u1.Name, u1.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			return mgr
		}},
		// 批量创建，每次创建一条
		{name: "create_2", args: args{[]*user{u1, u2}, 1}, wantErr: false, before: func() repo.Manager {
			mgr, mock := setupRepository()
			defer assert.NoError(t, mock.ExpectationsWereMet()) // 最后调用此方法，所有的期望值都匹配
			mock.ExpectBegin()
			//期望设置参数并成功执行sql，最后事务会提交
			mock.ExpectExec("INSERT INTO `users`").WithArgs(u1.CreatedAt, u1.UpdatedAt, u1.DeletedAt, u1.Name, u1.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("INSERT INTO `users`").WithArgs(u2.CreatedAt, u2.UpdatedAt, u2.DeletedAt, u2.Name, u2.ID).WillReturnResult(sqlmock.NewResult(2, 1))
			mock.ExpectCommit()
			return mgr
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := tt.before()
			if err := mgr.CreateBatchSize(tt.args.v, tt.args.size); (err != nil) != tt.wantErr {
				t.Errorf("CreateBatchSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_CreateBatchSizeTx(t *testing.T) {
	type args struct {
		v    any
		size int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		before  func() (repo.Manager, *gorm.DB)
	}{
		{name: "no_data", args: args{v: []*user{}, size: 0}, wantErr: false, before: func() (repo.Manager, *gorm.DB) {
			mgr, mock := setupRepository()
			defer assert.NoError(t, mock.ExpectationsWereMet())
			mock.ExpectBegin()    // 数据库调用mock开始
			mock.ExpectRollback() // 设置期望值，期望回滚。空slice执行会失败，数据库会回滚
			return mgr, mock.DB
		}},
		{name: "create_1", args: args{v: []*user{u1}, size: 1}, wantErr: false, before: func() (repo.Manager, *gorm.DB) {
			mgr, mock := setupRepository()
			defer assert.NoError(t, mock.ExpectationsWereMet())
			mock.ExpectBegin()
			// 期望设置参数并成功执行sql，最后事务会提交
			mock.ExpectExec("INSERT INTO `users`").WithArgs(u1.CreatedAt, u1.UpdatedAt, u1.DeletedAt, u1.Name, u1.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			return mgr, mock.DB
		}},
		// 批量创建，每次创建一条
		{name: "create_2", args: args{[]*user{u1, u2}, 1}, wantErr: false, before: func() (repo.Manager, *gorm.DB) {
			mgr, mock := setupRepository()
			defer assert.NoError(t, mock.ExpectationsWereMet()) // 最后调用此方法，所有的期望值都匹配
			mock.ExpectBegin()
			//期望设置参数并成功执行sql，最后事务会提交
			mock.ExpectExec("INSERT INTO `users`").WithArgs(u1.CreatedAt, u1.UpdatedAt, u1.DeletedAt, u1.Name, u1.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("INSERT INTO `users`").WithArgs(u2.CreatedAt, u2.UpdatedAt, u2.DeletedAt, u2.Name, u2.ID).WillReturnResult(sqlmock.NewResult(2, 1))
			mock.ExpectCommit()
			return mgr, mock.DB
		}},
	}
	for _, tt := range tests {
		t.Run("nil tx", func(t *testing.T) {
			mgr, _ := setupRepository()
			assert.PanicsWithError(t, "runtime error: invalid memory address or nil pointer dereference", func() {
				_ = mgr.CreateBatchSizeTx(nil, u1, 1)
			})
		})
		t.Run(tt.name, func(t *testing.T) {
			mgr, db := tt.before()
			if err := mgr.CreateBatchSizeTx(db, tt.args.v, tt.args.size); (err != nil) != tt.wantErr {
				t.Errorf("CreateBatchSizeTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_CreateOne(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
		before  func() repo.Manager
	}{
		{name: "create_1", args: args{v: u1}, want: u1, wantErr: false, before: func() repo.Manager {
			mgr, mock := setupRepository()
			defer assert.NoError(t, mock.ExpectationsWereMet()) // 最后调用此方法，所有的期望值都匹配

			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO `users`").WithArgs(u1.CreatedAt, u1.UpdatedAt, u1.DeletedAt, u1.Name, u1.ID).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			return mgr
		}},
	}
	for _, tt := range tests {
		t.Run("nil param", func(t *testing.T) {
			mgr, mock := setupRepository()
			defer assert.NoError(t, mock.ExpectationsWereMet()) // 最后调用此方法，所有的期望值都匹配

			_, err := mgr.CreateOne(nil)
			assert.Error(t, err)
		})
		t.Run(tt.name, func(t *testing.T) {
			mgr := tt.before()
			got, err := mgr.CreateOne(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_DeleteById(t *testing.T) {
	type args struct {
		v  any
		id uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		before  func() repo.Manager
	}{
		{name: "delete_1", args: args{v: u1, id: u1.ID}, wantErr: false, before: func() repo.Manager {
			mgr, mock := setupRepository()
			defer assert.NoError(t, mock.ExpectationsWereMet()) // 最后调用此方法，所有的期望值都匹配

			u1.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE `users` SET `deleted_at`=").WithArgs(test.AnyTime{}, u1.ID).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
			return mgr
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := tt.before()
			if err := mgr.DeleteById(tt.args.v, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_Find(t *testing.T) {
	columns := []string{"id", "created_at", "updated_at", "deleted_at", "name"}
	type args struct {
		v    any
		cond cond.Condition
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
		before  func() repo.Manager
	}{
		{name: "find_id_eq", args: args{v: u1, cond: cond.New("id = ?", u1.ID)}, want: []*user{u1}, wantErr: false, before: func() repo.Manager {
			mgr, mock := setupRepository()
			defer assert.NoError(t, mock.ExpectationsWereMet()) // 最后调用此方法，所有的期望值都匹配

			// mock 查询语句和参数
			mock.ExpectQuery("SELECT * FROM `users` WHERE id = \\? AND `users`.`deleted_at` IS NULL").WithArgs(u1.ID).
				// 返回的行数
				WillReturnRows(sqlmock.NewRows(columns).AddRow(u1.ID, u1.CreatedAt, u1.UpdatedAt, u1.DeletedAt, u1.Name))
			return mgr
		}},
		{name: "find_name_eq", args: args{v: u1, cond: cond.New("name = ?", u1.Name)}, want: u1, wantErr: false},
		{name: "find_name_like", args: args{v: u1, cond: cond.New("name like ?", u1.Name)}, want: u1, wantErr: false},
		{name: "find_name_like_2", args: args{v: u1, cond: cond.Wrap().Like("name", u1.Name)}, want: u1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := tt.before()
			got, err := mgr.Find(tt.args.v, tt.args.cond)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_FindOne(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v    any
		cond cond.Condition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.FindOne(tt.args.v, tt.args.cond)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_FindOneTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx   *gorm.DB
		v    any
		cond cond.Condition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.FindOneTx(tt.args.tx, tt.args.v, tt.args.cond)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOneTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOneTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_FindRaw(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v    any
		sql  string
		args []any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.FindRaw(tt.args.v, tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindRaw() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_FindRawTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx   *gorm.DB
		v    any
		sql  string
		args []any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.FindRawTx(tt.args.tx, tt.args.v, tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindRawTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindRawTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_FindTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx   *gorm.DB
		v    any
		cond cond.Condition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.FindTx(tt.args.tx, tt.args.v, tt.args.cond)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_GetById(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v  any
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.GetById(tt.args.v, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_GetByIdExact(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v  any
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.GetByIdExact(tt.args.v, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByIdExact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByIdExact() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_GetByIdExactTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		v  any
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.GetByIdExactTx(tt.args.tx, tt.args.v, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByIdExactTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByIdExactTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_GetByIdTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		v  any
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.GetByIdTx(tt.args.tx, tt.args.v, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByIdTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByIdTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_RemoveById(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v  any
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.RemoveById(tt.args.v, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("RemoveById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_RemoveByIdTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		v  any
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.RemoveByIdTx(tt.args.tx, tt.args.v, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("RemoveByIdTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_Save(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.Save(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Save() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_SaveTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx *gorm.DB
		v  any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.SaveTx(tt.args.tx, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_Tx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		f func(s repo.Manager, tx *gorm.DB) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.Tx(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("Tx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_UpdateField(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v     any
		field string
		val   any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.UpdateField(tt.args.v, tt.args.field, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("UpdateField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_UpdateFieldCond(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v     any
		cond  cond.Condition
		field string
		val   any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.UpdateFieldCond(tt.args.v, tt.args.cond, tt.args.field, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFieldCond() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_UpdateFieldCondTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx    *gorm.DB
		v     any
		cond  cond.Condition
		field string
		val   any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.UpdateFieldCondTx(tt.args.tx, tt.args.v, tt.args.cond, tt.args.field, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFieldCondTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_UpdateFieldTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx    *gorm.DB
		v     any
		field string
		val   any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.UpdateFieldTx(tt.args.tx, tt.args.v, tt.args.field, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFieldTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_UpdateFields(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v    any
		cond cond.Condition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.UpdateFields(tt.args.v, tt.args.cond); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_UpdateFieldsExact(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v      any
		fields map[string]any
		cond   cond.Condition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.UpdateFieldsExact(tt.args.v, tt.args.fields, tt.args.cond); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFieldsExact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_UpdateFieldsExactTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx     *gorm.DB
		v      any
		fields map[string]any
		cond   cond.Condition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.UpdateFieldsExactTx(tt.args.tx, tt.args.v, tt.args.fields, tt.args.cond); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFieldsExactTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_UpdateFieldsTx(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx   *gorm.DB
		v    any
		cond cond.Condition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if err := c.UpdateFieldsTx(tt.args.tx, tt.args.v, tt.args.cond); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFieldsTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGormManager_clone(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		tx *gorm.DB
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   repo.Manager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			if got := c.Clone(tt.args.tx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormManager_tableName(t *testing.T) {
	type fields struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &repo.GormManager{}
			got, err := c.TableName(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("tableName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tableName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGormManager(t *testing.T) {
	type args struct {
		rdb *gorm.DB
		wdb *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want repo.Manager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := repo.NewGormManager(tt.args.rdb, tt.args.wdb); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGormManager() = %v, want %v", got, tt.want)
			}
		})
	}
}
