package common

type Transfer struct {
	XferStart string
	XferTime  string
	XferBytes int
	XferIter  int
	Generator string
}

type Result struct {
	Protocol string
	Server   string
	Burst    int
	Start    string
	Times    []Transfer
}
