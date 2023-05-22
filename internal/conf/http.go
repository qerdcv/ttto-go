package conf

type HTTP struct {
	Addr string `required:"true" split_words:"true"`
}
