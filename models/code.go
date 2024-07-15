package models

const (
	NoRequestBody = iota + 5000
	MissingParameter
	CreateUserFailed
	LoginFailed
	BusinessFailed
	Web3VerifyFailed
	Web3LoginTimeout
	Web3ClinetIdError
)
