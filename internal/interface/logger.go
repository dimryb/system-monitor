package interfaces

//go:generate mockgen -source=logger.go -package=mocks -destination=../../mocks/mock_logger.go
type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
	Fatal(string, ...any)

	Debugf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
}
