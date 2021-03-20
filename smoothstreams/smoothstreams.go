package smoothstreams

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Site     string `yaml:"site"`
	Server   string `yaml:"server"`
}

type Client struct{}

func New() (*Client, error) {
	return nil, nil
}
