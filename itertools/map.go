package itertools

func FromMap[K comparable, V any](m map[K]V) Iterator[Pair[K, V]] {
	items := make([]Pair[K, V], 0, len(m))
	for key, value := range m {
		items = append(items, MakePair(key, value))
	}
	return FromSlice(items)
}

func ToMap[K comparable, V any](iter Iterator[Pair[K, V]]) map[K]V {
	m := make(map[K]V)
	for iter.Next() {
		pair := iter.Value()
		m[pair.first] = pair.second
	}
	return m
}
