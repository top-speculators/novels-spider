package interfaces

// 封装 Helper
// 所有用到辅助函数的地方都依赖于此接口，而不依赖于具体的包，如 utils
// 这样底层实现的包随时可替换

type Helper interface {
	SetLogConfig(t bool, path string)
	LoadConfig(path string) error
	GetConfig(s string) interface{}
}
