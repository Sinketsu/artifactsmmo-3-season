package strategy

type Strategy interface {
	Name() string
	Do() error
}
