package utils

import (
	"errors"
)

var (
	ErrUndefinedOS	    = errors.New("undefined operation system")
	ErrUndefinedCloud	= errors.New("undefined cloud provider")
	ErrIncorrectJob   	= errors.New("incorrect job")
	ErrShellInternal	= errors.New("internal shell error")
	ErrCloudInternal	= errors.New("internal cloud error")
	ErrInternal     	= errors.New("internal error")
	ErrForbidden        = errors.New("forbidden access")
	ErrBadToken         = errors.New("bad token")
	ErrGone             = errors.New("resource gone")
	ErrBusy             = errors.New("resource is busy")
)