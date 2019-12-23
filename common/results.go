package common

type Transfer struct {
	XferStart string
	XferTime  string
	XferBytes int
	XferIter  int
}

type Result struct {
	Protocol string
	Server   string
	Burst    int
	Start    string
	Times    []Transfer
}
