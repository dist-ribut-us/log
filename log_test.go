package log

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	gl "log"
	"testing"
)

func TestLog(t *testing.T) {
	b := &bytes.Buffer{}
	l := New(b)
	assert.NotNil(t, l)
	l = l.Child("test")
	assert.NotNil(t, l)
	assert.Equal(t, `"test"`, l.data)
	l.Info("foo", "bar")
	l.Info(KV{"pi", 3.1415})
	s := b.String()
	assert.Contains(t, s, ` "test" "foo" "bar"`)
	assert.Contains(t, s, ` pi=3.1415`)
	assert.NotContains(t, s, "  ")
}

func TestTwoChildren(t *testing.T) {
	b := &bytes.Buffer{}
	l := New(b)
	assert.NotNil(t, l)
	l = l.Child("test")
	l1 := l.Child("child1")
	l2 := l.Child("child2")

	l1.Info("test1")
	s := b.String()
	assert.Contains(t, s, "child1")
	assert.NotContains(t, s, "child2")
	b.Reset()

	l2.Info("test2")
	s = b.String()
	assert.Contains(t, s, "child2")
	assert.NotContains(t, s, "child1")
	b.Reset()
}

func TestNil(t *testing.T) {
	var l *Log
	assert.Nil(t, l.Child())
	l.Info("test")
	l.Close()
}

func TestChangeLogForChild(t *testing.T) {
	b1 := &bytes.Buffer{}
	p := New(b1)
	c := p.Child("child")
	b2 := &bytes.Buffer{}
	p.To(b2)
	c.Info("test")

	s := b2.String()
	assert.Contains(t, s, "test")
	assert.Contains(t, s, "child")

	p.Info("foo")
	s = b2.String()
	assert.Contains(t, s, "foo")
}

func TestDebug(t *testing.T) {
	var p *Log
	assert.False(t, p.GetDebug())

	b := &bytes.Buffer{}
	p = New(b)
	assert.False(t, p.GetDebug())

	p.Debug("test1")
	s := b.String()
	assert.Equal(t, "", s)

	p.SetDebug(true)
	p.Debug("test2")
	s = b.String()
	assert.Contains(t, s, "test2")
}

func TestGoLog(t *testing.T) {
	b := &bytes.Buffer{}
	To(b)

	s := b.String()
	assert.Equal(t, "", s)

	Go()
	gl.Print("test2")
	s = b.String()
	assert.Contains(t, s, "test2")
}

func TestTrace(t *testing.T) {
	b := &bytes.Buffer{}
	l := New(b)
	l.SetTrace(0, 1, 1)
	A(l)
	s := b.String()
	assert.Contains(t, s, "INFO log/log_test.go:")
	assert.Contains(t, s, "log.A")

	b.Reset()
	l.SetTrace(0, 2, 1)
	B(l)
	s = b.String()
	assert.Contains(t, s, "\n\tlog/log_test.go:")
	assert.Contains(t, s, "log.A")
	assert.Contains(t, s, "log.B")
}

func A(lgr *Log) {
	lgr.Info()
}

func B(lgr *Log) {
	A(lgr)
}
