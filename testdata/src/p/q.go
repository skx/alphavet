package p

type Steve string

func (s Steve) Foo() {
}
func (s Steve) Bar() { // want "function Bar on receiver s should have been before Foo"
}
