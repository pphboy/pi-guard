package sys

type ConfigTool interface {
	GetConfig() SysConfig
}

func NewConfigTool() ConfigTool {
	return nil
}

type SysConfig struct {
}
