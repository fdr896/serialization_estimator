package support

type StringSet struct {
	index map[string]struct{}
}

func NewStringSet(strings []string) *StringSet {
	ss := StringSet{
		index: make(map[string]struct{}),
	}
	for _, str := range strings {
		ss.index[str] = struct{}{}
	}
	
	return &ss
}

func (ss *StringSet) Contains(str string) bool {
	_, ok := ss.index[str]
	return ok
}

func (ss *StringSet) Keys() []string {
	keys := make([]string, 0)
	for k := range ss.index {
		keys = append(keys, k)
	}

	return keys
}
