package protocol

type Option func(m *Message)

func WithElements(e ...*Element) Option {
	return func(m *Message) {
		m.Elements = append(m.Elements, e...)
	}
}

func WithProtocolData(strSlice ...string) Option {
	return func(m *Message) {
		for _, s := range strSlice {
			m.originalData.WriteString(s)
		}
	}
}
