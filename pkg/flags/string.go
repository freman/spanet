package flags

type String struct {
	set   bool
	value string
}

func (s *String) Set(x string) error {
	s.value = x
	s.set = true
	return nil
}

func (s *String) String() string {
	return s.value
}

func (s *String) IsSet() bool {
	return s.set
}
