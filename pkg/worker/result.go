package worker

type Result interface {
	PrintRecreate()
	Sha1() string
	Object() []byte
}

type Worker interface {
	Count() uint64
	Work(chan Result)
}
