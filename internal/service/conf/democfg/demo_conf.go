package democfg

import (
	"goboot/internal/service/conf"
	"goboot/pkg/printer"
	"strconv"
)

var DemoCfg = &demoCfg{}

type demoCfg struct {
}

func (c *demoCfg) UsernameLen() (int, int) {
	return minUsernameLen, maxUsernameLen
}

// ==============================
// 配置示例
// ==============================

const (
	minUsernameLen = 6
	maxUsernameLen = 20
)

var (
	usernameLenCfg = cfg{}
)

func init() {
	conf.RegisterConfigurator(&configurator{})
}

type configurator struct {
}

func (*configurator) Title() string {
	return "用户名长度配置示例"
}

type cfg struct {
	min int
	max int
}

func (c *configurator) Configure(print bool) {
	usernameLenCfg.min = minUsernameLen
	usernameLenCfg.max = maxUsernameLen
	if print {
		c.print()
	}
}

func (c *configurator) print() {
	printer.NewLine()
	printer.NewSepLine()
	printer.Println("表: %s", c.Title())
	w := 10
	printer.Printwln(w, "MIN", "MAX")
	printer.Printwln(w, strconv.Itoa(usernameLenCfg.min), strconv.Itoa(usernameLenCfg.max))
	printer.NewSepLine()
}
