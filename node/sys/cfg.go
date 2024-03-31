package sys

type ConfigTool interface {
	GetConfig() SysConfig
}

func NewConfigTool() ConfigTool {
	return nil
}

// ! FIXME: 等以后有空再将这个相关变量改成配置文件形式的
type SysConfig struct {
	RootDomain string
}
