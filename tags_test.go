package errtags

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestColors(t *testing.T) {
	var redAlertClass = NewTag("WARNING, WARNING, RED ALERT")
	var blueAlertClass = NewTag("nothing to worry about")

	redError := redAlertClass.Tag(errors.New("enemy spotted"))
	blueError := blueAlertClass.Tag(errors.New("mark has slipped on a banana peel"))
	nonClassError := errors.New("some error")

	assert.ErrorIs(t, redError, redAlertClass)
	assert.ErrorIs(t, blueError, blueAlertClass)
	assert.NotErrorIs(t, redError, blueAlertClass)
	assert.NotErrorIs(t, blueError, redAlertClass)
	assert.NotErrorIs(t, redError, nonClassError)
	assert.NotErrorIs(t, blueError, nonClassError)
}

func TestSameMessage(t *testing.T) {
	var sameMessageClassA = NewTag("same message")
	var sameMessageClassB = NewTag("same message")

	sameMessageErrorA := sameMessageClassA.Tag(errors.New("same message"))
	sameMessageErrorB := sameMessageClassB.Tag(errors.New("same message"))
	nonClassError := errors.New("same message")

	assert.ErrorIs(t, sameMessageErrorA, sameMessageClassA)
	assert.ErrorIs(t, sameMessageErrorB, sameMessageClassB)
	assert.NotErrorIs(t, sameMessageErrorA, sameMessageClassB)
	assert.NotErrorIs(t, sameMessageErrorB, sameMessageClassA)
	assert.NotErrorIs(t, sameMessageErrorA, nonClassError)
	assert.NotErrorIs(t, sameMessageErrorB, nonClassError)
}

func TestNoMessage(t *testing.T) {
	var noMessageClassA = NewTag()
	var noMessageClassB = NewTag()

	noMessageErrorA := noMessageClassA.Tag(errors.New("some message"))
	noMessageErrorB := noMessageClassB.Tag(errors.New("some message"))

	assert.ErrorIs(t, noMessageErrorA, noMessageClassA)
	assert.ErrorIs(t, noMessageErrorB, noMessageClassB)
	assert.NotErrorIs(t, noMessageErrorA, noMessageClassB)
	assert.NotErrorIs(t, noMessageErrorB, noMessageClassA)
}

func TestUnwrapped(t *testing.T) {
	var someClassA = NewTag("some message")
	var someClassB = NewTag("some message")

	assert.NotErrorIs(t, someClassA, someClassB)
}

func TestStack(t *testing.T) {
	var someClass = NewTag()

	baseError := errors.New("some message")
	wrappedError := someClass.Tag(baseError)

	baseErrorPrint := fmt.Sprintf("%+v", baseError)
	wrappedErrorPrint := fmt.Sprintf("%+v", wrappedError)

	// remove first lines, need to compare only the stack
	_, baseErrorPrint, _ = strings.Cut(baseErrorPrint, "\n")
	_, wrappedErrorPrint, _ = strings.Cut(wrappedErrorPrint, "\n")

	assert.Equal(t, baseErrorPrint, wrappedErrorPrint)
}

func TestMessage(t *testing.T) {
	var someClass = NewTag("class message")

	baseError := errors.New("some message")
	wrappedError := someClass.Tag(baseError)

	assert.Equal(t, "class message: some message", fmt.Sprintf("%v", wrappedError))
}

func TestMessageNoClass(t *testing.T) {
	var someClass = NewTag()

	baseError := errors.New("some message")
	wrappedError := someClass.Tag(baseError)

	// error classes without additional message should create identical formatted errors
	// as the base error
	assert.Equal(t, fmt.Sprintf("%+v", baseError), fmt.Sprintf("%+v", wrappedError))
}

func TestWrappingNil(t *testing.T) {
	var someClass = NewTag("class message")

	wrappedError := someClass.Tag(nil)

	assert.Nil(t, wrappedError)

	wrappedError = someClass.TagWithMessage(nil, "new message")

	assert.Nil(t, wrappedError)
}

func TestWithMessage(t *testing.T) {
	var someClass = NewTag("class message")

	baseError := errors.New("some message")
	wrappedError := someClass.TagWithMessage(baseError, "new message")

	assert.Equal(t, "new message: some message", fmt.Sprintf("%v", wrappedError))
	assert.ErrorIs(t, wrappedError, someClass)
}

func TestWithMessageNoClassMessage(t *testing.T) {
	var someClass = NewTag()

	baseError := errors.New("some message")
	wrappedError := someClass.TagWithMessage(baseError, "new message")

	assert.Equal(t, "new message: some message", fmt.Sprintf("%v", wrappedError))
	assert.ErrorIs(t, wrappedError, someClass)
}

func TestCauseUnwrap(t *testing.T) {
	var someClass = NewTag("class message")

	baseError := errors.New("some message")
	wrappedError := someClass.Tag(baseError)

	assert.Equal(t, baseError, errors.Cause(wrappedError))
	assert.Equal(t, baseError, errors.Unwrap(wrappedError))
}

func TestError(t *testing.T) {
	var someClass = NewTag("class message")

	assert.Equal(t, "class message", someClass.Error())

	baseError := errors.New("some message")
	wrappedError := someClass.Tag(baseError)

	assert.Equal(t, "class message: some message", wrappedError.Error())
}

func TestWithTags(t *testing.T) {
	redTag := NewTag("red tag")
	blueTag := NewTag("blue tag")

	someError := errors.New("some error")

	colorfulTagged := WithTags(someError, redTag, blueTag)

	assert.ErrorIs(t, colorfulTagged, redTag)
	assert.ErrorIs(t, colorfulTagged, blueTag)
}

func TestWithTagsNested(t *testing.T) {
	redTag := NewTag("red tag")
	blueTag := NewTag("blue tag")

	someError := errors.New("some error")

	withRedTag := WithTags(someError, redTag)
	withBlueAndRedTag := WithTags(withRedTag, blueTag)

	assert.ErrorIs(t, withBlueAndRedTag, someError)
	assert.ErrorIs(t, withBlueAndRedTag, redTag)
	assert.ErrorIs(t, withBlueAndRedTag, blueTag)

	assert.ErrorIs(t, withRedTag, redTag)
	assert.ErrorIs(t, withRedTag, someError)

	assert.Equal(t, "red tag: some error", withRedTag.Error())

	assert.NotErrorIs(t, withRedTag, blueTag)

	assert.Equal(t, "blue tag: red tag: some error", withBlueAndRedTag.Error())
}

func TestTags(t *testing.T) {
	redTag := NewTag("red tag")
	blueTag := NewTag("blue tag")

	someError := errors.New("some error")

	withTags := WithTags(someError, redTag, blueTag)

	assert.Equal(t, []*Tag{redTag, blueTag}, withTags.(*Tag).Tags())
	assert.ErrorIs(t, withTags, redTag)
	assert.ErrorIs(t, withTags, blueTag)
	assert.Equal(t, "red tag: blue tag: some error", withTags.Error())
}

func TestWithTagsAndMessage(t *testing.T) {
	redTag := NewTag("red tag")
	blueTag := NewTag("blue tag")

	someError := errors.New("some error")

	colorfulTaggedWithMessage := WithTagsAndMessage(someError, "message", redTag, blueTag)

	assert.Equal(t, []*Tag{redTag, blueTag}, colorfulTaggedWithMessage.(*Tag).Tags())
	assert.ErrorIs(t, colorfulTaggedWithMessage, redTag)
	assert.ErrorIs(t, colorfulTaggedWithMessage, blueTag)
	assert.Equal(t, "message: some error", colorfulTaggedWithMessage.Error())
}

func TestUnion(t *testing.T) {
	redTag := NewTag("red tag")
	blueTag := NewTag("blue tag")

	someError := errors.New("some error")

	colorfulTag := UnionTag(redTag, blueTag)

	colorfulTagged := colorfulTag.Tag(someError)
	redTagged := redTag.Tag(someError)
	blueTagged := blueTag.Tag(someError)

	assert.Equal(t, []*Tag{redTag, blueTag}, colorfulTag.Tags())
	assert.ErrorIs(t, colorfulTag, redTag)
	assert.ErrorIs(t, colorfulTag, blueTag)
	assert.ErrorIs(t, colorfulTagged, redTag)
	assert.ErrorIs(t, colorfulTagged, blueTag)
	assert.ErrorIs(t, colorfulTagged, colorfulTag)

	assert.NotErrorIs(t, redTag, colorfulTag)
	assert.NotErrorIs(t, redTagged, colorfulTag)
	assert.NotErrorIs(t, blueTag, colorfulTagged)
	assert.NotErrorIs(t, blueTagged, colorfulTag)

	colorfulTaggedWithMessage := UnionTag(redTag, blueTag).TagWithMessage(someError, "message")

	assert.Equal(t, []*Tag{redTag, blueTag}, colorfulTaggedWithMessage.(*Tag).Tags())
	assert.ErrorIs(t, colorfulTaggedWithMessage, redTag)
	assert.ErrorIs(t, colorfulTaggedWithMessage, blueTag)
	assert.ErrorIs(t, colorfulTaggedWithMessage, colorfulTag)

	assert.NotErrorIs(t, redTag, colorfulTaggedWithMessage)
	assert.NotErrorIs(t, blueTag, colorfulTaggedWithMessage)
	assert.NotErrorIs(t, redTagged, colorfulTaggedWithMessage)
	assert.NotErrorIs(t, blueTagged, colorfulTaggedWithMessage)
	assert.Equal(t, "message: some error", colorfulTaggedWithMessage.Error())
}
