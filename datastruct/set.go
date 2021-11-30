package datastruct

type dictEntry struct {
	value interface{}
}

func (de *dictEntry) getValue() interface{} {
	return de.value
}

func (de *dictEntry) setValue(v interface{}) {
	de.value = v
}

func withValue(v interface{}) *dictEntry {
	return &dictEntry{value: v}
}
