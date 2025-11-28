package repository

type M2MSyncPolicy int

const (
	SyncAlways M2MSyncPolicy = iota
	SyncIfProvided
)
