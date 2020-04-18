package options

type Target int

const (
	ES5 Target = iota
	ES2015
	ES2016
	ES2017
	ES2018
	ES2019
	ES2020
)

type Options struct {
	Target Target
	Minify bool
}
