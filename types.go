package main

type Project struct {
	timeUnit TimeUnit
	currency Currency
	team     map[string]Resource
	risks    map[string]float32
	tasks    []Task
}

type TimeUnit uint8

const (
	hr TimeUnit = iota
	day
)

type Currency uint8

const (
	usd Currency = iota
	eur
)

type Resource struct {
	Title   string
	Rate    float64
	Count   uint8
	Formula string
}

type Task struct {
	Category string
	Title    string
	Risk     string
	Work     map[string]float32 // resource -> time units
}
