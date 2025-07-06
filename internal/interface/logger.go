package interfaces

//go:generate mockgen -source=logger.go -package=mocks -destination=../../mocks/mock_logger.go
type Logger interface {
	Debugf(string, ...any)
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
}
