package itertools

type rangeConfig struct {
	start int
	stop  int
	step  int
}

func (config rangeConfig) Unpack() (start, stop, step int) {
	return config.start, config.stop, config.step
}

type RangeOption func(*rangeConfig)

func getConfig(options []RangeOption) rangeConfig {
	config := rangeConfig{start: 0, stop: 0, step: 1}
	for _, option := range options {
		option(&config)
	}
	return config
}

func Start(start int) RangeOption {
	return func(config *rangeConfig) {
		config.start = start
	}
}

func Stop(start int) RangeOption {
	return func(config *rangeConfig) {
		config.stop = start
	}
}

func Step(start int) RangeOption {
	return func(config *rangeConfig) {
		config.step = start
	}
}

func Range(options ...RangeOption) Iterator[int] {
	config := getConfig(options)
	start, stop, step := config.start, config.stop, config.step

	value := start
	var sign int
	if step < 0 && start > stop {
		sign = -1
	} else if step > 0 && start < stop {
		sign = 1
	} else {
		return EmptyIterator[int]()
	}
	advance := func() (bool, int) {
		if sign*value >= sign*stop {
			return false, 0
		} else {
			retValue := value
			value += step
			return true, retValue
		}
	}
	return FromAdvance(advance)
}
