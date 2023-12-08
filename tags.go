package errtags

import (
	"fmt"
	"io"
	"strings"
)

type Tag struct {
	causer      error
	msg         string
	tags        []*Tag
	msgOverride bool
}

// Tag adds tags to the error.
// Stack trace of the error is preserved.
// If tag has no message, the original error message is also unchanged.
// If tag has a message, the original error message is appended to it.
func (e *Tag) Tag(err error) error {
	if err == nil {
		return nil
	}

	return &Tag{
		causer: err,
		msg:    e.msg,
		tags:   e.tags,
	}
}

func (e *Tag) message() string {
	if e.msgOverride {
		return e.msg
	}

	builder := strings.Builder{}
	ln := len(e.tags)
	for i, tag := range e.tags {
		if tag.msg != "" {
			_, _ = builder.WriteString(tag.msg)
			if i < ln-1 {
				_, _ = builder.WriteString(": ")
			}
		}
	}
	return builder.String()
}

// base error interface

func (e *Tag) Error() string {
	ownMessage := e.message()

	if e.causer == nil {
		return ownMessage
	}

	return ownMessage + ": " + e.causer.Error()
}

// formatting interface

func (e *Tag) Format(s fmt.State, verb rune) {
	ownMessage := e.message()
	if ownMessage != "" {
		_, _ = io.WriteString(s, ownMessage)
		_, _ = io.WriteString(s, ": ")
	}

	if fmtr, ok := e.causer.(fmt.Formatter); ok {
		fmtr.Format(s, verb)
		return
	}
}

// interface for error libraries' equality checks

func (e *Tag) Is(err error) bool {
	//goland:noinspection GoTypeAssertionOnErrors
	if other, ok := err.(*Tag); ok {
		// true if either wc.tags or e.tags is a subset of the other
		return isSubset(e.tags, other.tags)
	}

	return false
}

// TagWithMessage tags the error and adds a message.
func (e *Tag) TagWithMessage(err error, msg string) error {
	if err == nil {
		return nil
	}

	tag := Tag{
		causer:      err,
		msg:         msg,
		msgOverride: true,
		tags:        e.tags,
	}

	return &tag
}

func NewTag(msg ...string) *Tag {
	tag := Tag{}
	if len(msg) > 0 {
		tag.msg = strings.Join(msg, ": ")
	}
	tag.tags = []*Tag{&tag}
	return &tag
}

func (e *Tag) Tags() []*Tag {
	return e.tags
}

func WithTags(err error, tags ...*Tag) error {
	if err == nil {
		return nil
	}

	return &Tag{
		causer: err,
		tags:   tags,
	}
}

func WithTagsAndMessage(err error, msg string, tags ...*Tag) error {
	if err == nil {
		return nil
	}

	//goland:noinspection GoTypeAssertionOnErrors
	tagged := WithTags(err, tags...).(*Tag)
	tagged.msg = msg
	tagged.msgOverride = true

	return tagged
}

func UnionTag(tags ...*Tag) *Tag {
	return &Tag{
		tags: tags,
	}
}

// Satisfy common interfaces

func (e *Tag) Cause() error {
	return e.causer
}

func (e *Tag) Unwrap() error {
	return e.causer
}
