package conf

type Auth struct {
	Secret []byte `required:"true"`
}
