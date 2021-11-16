package datastruct

type Operator interface {
	Save(key, value string, options ...interface{})
	Get(key string) (value string, err error)
	Increase(key string, value float64)
}
