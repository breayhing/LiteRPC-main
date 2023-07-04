package main

type ARGV struct {
	ip         string
	port       int
	methodType string
	nodeNum    int
}

// ARGV默认值
var argv = ARGV{
	ip:         "",
	port:       0,
	methodType: "string",
	nodeNum:    1,
}

// func (*ARGV) Option()
