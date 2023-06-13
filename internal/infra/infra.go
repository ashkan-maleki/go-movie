package infra

type Infra struct {
	Config *config
}

func New() (*Infra, error) {
	conf, err := newConfig()
	if err != nil {
		return nil, err
	}
	return &Infra{
		Config: conf,
	}, nil
}
