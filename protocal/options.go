package protocal

type Option func(m *Message)

func SimpleString(str string) Option {
	return func(m *Message) {
		m.Elements = append(m.Elements, NewSimpleStringElement(str))
	}
}

func BulkString(str string) Option {
	return func(m *Message) {
		m.Elements = append(m.Elements, NewBulkStringElement(str))
	}
}

func ErrorString(str string) Option {
	return func(m *Message) {
		m.Elements = append(m.Elements, NewErrorElement(str))
	}
}

func Integer(i int) Option {
	return func(m *Message) {
		m.Elements = append(m.Elements, NewIntegerElement(i))
	}
}

func Array(opts ...Option) Option {
	if len(opts) == 0 {
		return func(_ *Message) {}
	}

	return func(m *Message) {
		for _, opt := range opts {
			opt(m)
		}
	}
}