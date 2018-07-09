package loader

type loaderKey string

func (k loaderKey) String() string {
	return string(k)
}

const (
	UserLoaderKey loaderKey = "user"
)

// 需要对应proto文件中的node type
type GlobalType string

const (
	User GlobalType = "user"
)

type DataKey struct {
	string string
	raw    interface{}
}

func NewDataKey(key string, raw interface{}) *DataKey {
	return &DataKey{key, raw}
}

func (d *DataKey) String() string {
	return d.string
}

func (d *DataKey) Raw() interface{} {
	if d != nil {
		return d.raw
	}
	return nil
}
