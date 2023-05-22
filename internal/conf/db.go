package conf

import "net/url"

type DBConfig struct {
	Username string `required:"true" split_words:"true"`
	Password string `required:"true" split_words:"true"`
	DBName   string `required:"true" split_words:"true"`
	HOST     string `required:"true" split_words:"true"`
}

func (c *DBConfig) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.Username, c.Password),
		Host:   c.HOST,
		Path:   c.DBName,
	}

	u.RawQuery = "sslmode=disable"
	return u.String()
}
