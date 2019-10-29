package matcher

type Matcher interface {
	Match([]byte) bool
}
