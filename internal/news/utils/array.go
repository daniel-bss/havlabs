package utils

func Map[A, B any](this []A, fn func(b A) B) (result []B) {
	for _, v := range this {
		result = append(result, fn(v))
	}
	return
}
