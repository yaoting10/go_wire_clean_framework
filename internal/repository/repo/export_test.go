package repo

import (
	"fmt"
	"goboot/pkg/web3"
	"gorm.io/gorm"
	"testing"
)

func (c *GormManager) Clone(tx *gorm.DB) Manager {
	return c.clone(tx)
}

func (c *GormManager) TableName(v any) (string, error) {
	return c.tableName(v)
}

func Test(t *testing.T) {
	//w, _ := web3.GenerateWithMnemonic("fun spin teach candy goat struggle grid bus slender scale develop suit", "123456")
	w, _ := web3.GenerateWithMnemonic("shock idle level flight into little consider ill more brisk slow slide", "")
	//w, _ := web3.GenerateWithMnemonic("", "12345678")
	//w, _ := web3.GenerateWithMnemonic("", "12345678")
	//"0x32D879f48e17bF9eE07474d93b04670C1186b8B2"
	//"0x32D879f48e17bF9eE07474d93b04670C1186b8B2"
	fmt.Println(w)
}
