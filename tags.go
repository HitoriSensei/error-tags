package errtags

import (
	"fmt"
	"io"
	"reflect"
	"slices"
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

// Include extends the Tag with additional tags.
func (e *Tag) Include(tags ...*Tag) *Tag {
	e.tags = sort(append(e.tags, tags...))
	return e
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
		// true if other.tags is a sub-slice of e.tags
		return isSubSlice(e.tags, other.tags)
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

// NewTag creates a new Tag with optional message.
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

// WithTags tags the error with the given tags.
func WithTags(err error, tags ...*Tag) error {
	if err == nil {
		return nil
	}

	return &Tag{
		causer: err,
		tags:   sort(tags),
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

func UnionTag(tag *Tag, tags ...*Tag) *Tag {
	return &Tag{
		tags: sort(append(tags, tag)),
	}
}

// recursively extract tags from tags
func getAllTags(tags []*Tag, seen *map[*Tag]struct{}) []*Tag {
	allTags := make([]*Tag, 0, len(tags))

	if seen == nil {
		m := make(map[*Tag]struct{}, len(tags))
		seen = &m
	}

	for _, tag := range tags {
		if _, ok := (*seen)[tag]; ok {
			continue
		}
		(*seen)[tag] = struct{}{}
		allTags = append(allTags, tag)
		allTags = append(allTags, getAllTags(tag.tags, seen)...)
	}

	return allTags
}

func sort(tags []*Tag) []*Tag {
	// recursively extract tags from tags and remove duplicates
	allTags := uniq(getAllTags(tags, nil))

	// sort tags by their pointer value
	slices.SortFunc(allTags, func(i, j *Tag) int {
		return int(reflect.ValueOf(i).Pointer()) - int(reflect.ValueOf(j).Pointer())
	})

	return allTags
}

func uniq(tags []*Tag) []*Tag {
	seen := make(map[*Tag]struct{}, len(tags))
	unique := make([]*Tag, 0, len(tags))

	for _, tag := range tags {
		if _, ok := seen[tag]; !ok {
			seen[tag] = struct{}{}
			unique = append(unique, tag)
		}
	}

	return unique
}

// Satisfy common interfaces

func (e *Tag) Cause() error {
	return e.causer
}

func (e *Tag) Unwrap() error {
	return e.causer
}

func Equal(a, b *Tag) bool {
	return a.Is(b) && b.Is(a)
}
