package constant

import "errors"

var (
	ConnNonDial     error = errors.New("Conn haven't be dialed !")
	ConnHasBeDialed error = errors.New("Conn has be dialed !")
	TransConfNil    error = errors.New("Trans without config !")
)
