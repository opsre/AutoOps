package agent

// Agent版本信息
const (
	AgentVersion    = "2.0.0" // Agent版本号
	AgentVersionTag = "v2"    // Agent版本标签
	AgentBuildTime  = ""      // 构建时间(编译时可注入)
)

// Config 结构体定义
type Config struct {
	Processes []struct {
		Name string `yaml:"name"`
	} `yaml:"processes"`
	Endpoints []struct {
		URL    string `yaml:"url"`
		Method string `yaml:"method"`
	} `yaml:"endpoints"`
	Pushgateway struct {
		URL          string `yaml:"url"`
		PushInterval string `yaml:"push_interval"`
		Job          string `yaml:"job"`
		Instance     string `yaml:"instance"`
	} `yaml:"pushgateway"`
	Heartbeat struct {
		ServerURL string `yaml:"server_url"`
		Interval  string `yaml:"interval"`
		Token     string `yaml:"token"`
	} `yaml:"heartbeat"`
	CollectInterval string `yaml:"collect_interval"`
}

// 默认配置
var defaultConfig = Config{
	Processes: []struct {
		Name string `yaml:"name"`
	}{
		{Name: "nginx"},
		{Name: "mysql"},
		{Name: "redis"},
		{Name: "go"},
		{Name: "bash"},
	},
	Endpoints: []struct {
		URL    string `yaml:"url"`
		Method string `yaml:"method"`
	}{
		{URL: "http://localhost:9100/metrics", Method: "GET"},
	},
	Pushgateway: struct {
		URL          string `yaml:"url"`
		PushInterval string `yaml:"push_interval"`
		Job          string `yaml:"job"`
		Instance     string `yaml:"instance"`
	}{
		URL:          "", // 从配置文件或环境变量读取
		PushInterval: "30s",
		Job:          "devops-agent-%s",
		Instance:     "$HOSTNAME",
	},
	Heartbeat: struct {
		ServerURL string `yaml:"server_url"`
		Interval  string `yaml:"interval"`
		Token     string `yaml:"token"`
	}{
		ServerURL: "", // 从配置文件或环境变量读取
		Interval:  "30s",
		Token:     "agent-heartbeat-token-2024", // 默认token
	},
	CollectInterval: "15s",
}
