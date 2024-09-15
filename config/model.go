package config

type Config struct {
	SplitStrategy []SplitStrategy `yaml:"splitStrategy"`
}

type SplitStrategy struct {
	HostPattern string `yaml:"hostPattern"` // host匹配规则
	Type        int    `yaml:"type"`        // 1cookie 2header
	Field       string `yaml:"field"`       // 字段
}
