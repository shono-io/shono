package utils

import (
	"errors"
	"fmt"
)

var ErrTimeout = errors.New("operation timed out")

func NotFoundError(reason string, args ...any) error {
	return &notFoundError{fmt.Sprintf(reason, args...)}
}

type notFoundError struct {
	reason string
}

func (e *notFoundError) Error() string {
	return fmt.Sprintf("[NOT FOUND] %s", e.reason)
}

func InternalError(err error) error {
	return &internalError{err}
}

type internalError struct {
	err error
}

func (e *internalError) Error() string {
	return fmt.Sprintf("[INTERNAL] %s", e.err.Error())
}

func IllegalStateError(reason string, args ...any) error {
	return &illegalStateError{fmt.Sprintf(reason, args...)}
}

type illegalStateError struct {
	reason string
}

func (e *illegalStateError) Error() string {
	return fmt.Sprintf("[ILLEGAL STATE] %s", e.reason)
}

func NewFunctionalError(err error) *FunctionalError {
	return &FunctionalError{err, false}
}

/**
 * panic - if no other records may be processed due to the order of processing needed to be preserved
 * ack - ack the message and continue processing
 * no-ack / redeliver - do not ack the message, causing the message to be redelivered. This should also cause the
 * remaining messages in the batch to be redelivered.

 * connection to an external resource is lost. [blocking, technical, retryable]
 * message is malformed. [blocking, functional, non-retryable] !! POISON PILL !!
 * message can never be processed by external resource. [blocking, functional, non-retryable] !! POISON PILL !!
 */

/*
FunctionalError is an error that is not related to the system itself, but to logical processing

A FunctionalError can be blocking or not. A blocking error will prevent the system from continuing
processing messages.
*/
type FunctionalError struct {
	err      error
	blocking bool
}

func (f *FunctionalError) Error() string {
	if f.blocking {
		return fmt.Sprintf("[BLOCKING] %s", f.err.Error())
	}

	return f.err.Error()
}

func NewTechnicalError(err error) *TechnicalError {
	return &TechnicalError{err, false}
}

/*
TechnicalError is an error that is related to the system itself. It is a wrapper around an error
caused by a technical issue.

A TechnicalError can be blocking or not. A blocking technical error is also known as a poison pill
*/
type TechnicalError struct {
	err      error
	blocking bool
}

func (e *TechnicalError) Error() string {
	if e.blocking {
		return fmt.Sprintf("[BLOCKING] %s", e.err.Error())
	}

	return e.err.Error()
}

//func PoisonPill(evt *event.NewEvent, err error) error {
//	return &PoisonPillError{evt, err}
//}
//
//type PoisonPillError struct {
//	evt *event.NewEvent
//	err error
//}
//
//func (e *PoisonPillError) Error() string {
//	return fmt.Sprintf("[POISON PILL] %s", e.err.Error())
//}
