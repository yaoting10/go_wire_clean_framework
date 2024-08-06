package conf

var configurators []Configurator

func Configure() {
	printDetail := true
	//fmt.Println("正在初始化核心配置项...")
	for _, c := range configurators {
		//fmt.Printf("正在配置 [%s] ...\n", c.Title())
		c.Configure(printDetail)
	}
}

func RegisterConfigurator(c Configurator) {
	configurators = append(configurators, c)
}

type Configurator interface {
	Title() string
	Configure(print bool)
}

type Scope[T comparable] struct {
	Min T
	Max T
}

func NewScope[T comparable](min, max T) Scope[T] {
	return Scope[T]{Min: min, Max: max}
}
