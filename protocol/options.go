package protocol

type Option func(m *Message)

func WithElements(e ...*Element) Option {
	return func(m *Message) {
		m.Elements = append(m.Elements, e...)
	}
}

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

func Integer(is string) Option {
	return func(m *Message) {
		m.Elements = append(m.Elements, NewIntegerElement(is))
	}
}

func Array(opts ...Option) Option {
	if len(opts) == 0 {
		return func(_ *Message) {}
	}

	return func(m *Message) {
		m.Elements = append(m.Elements, NewArrayElement(len(opts)))
		for _, opt := range opts {
			opt(m)
		}
	}
}
