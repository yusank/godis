package protocol

type option func(m *Message)

func withElements(e ...*Element) option {
	return func(m *Message) {
		m.Elements = append(m.Elements, e...)
	}
}
