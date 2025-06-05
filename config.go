package es

type Config struct {
	Name      string
	Addresses []string
	Username  string
	Password  string
	Tracing   bool
	Metrics   bool
}
