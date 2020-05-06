package main

import "errors"

func fib(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("n cannot be less than 0")
	}

	if n < 2 {
		return n, nil
	}

	a, err := fib(n - 2)
	if err != nil {
		return 0, err
	}

	b, err := fib(n - 1)
	if err != nil {
		return 0, err
	}

	return a + b, nil
}
