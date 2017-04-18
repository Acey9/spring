package main

type Control struct {
	buf string
}

func NewCtrl(str string) (*Control, error) {
	ctrl := &Control{str}
	return ctrl, nil
}
