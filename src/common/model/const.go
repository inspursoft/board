package model

const (
	ProjectAdmin = int64(iota + 1)
	Developer
	Visitor
	ServiceStart = int64(iota + 1)
	ServiceStop
)
