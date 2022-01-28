package itertools

func FromSet[K comparable](m map[K]struct{}) Iterator[K] {
	items := make([]K, 0, len(m))
	for key, _ := range m {
		items = append(items, key)
	}
	return FromSlice(items)
}

func ToSet[K comparable](iter Iterator[K]) map[K]struct{} {
	m := make(map[K]struct{})
	for iter.Next() {
		k := iter.Value()
		m[k] = struct{}{}
	}
	return m
}
