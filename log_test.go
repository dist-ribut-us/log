package log

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	gl "log"
	"os"
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

func TestFile(t *testing.T) {
	l, err := File("test.log")
	assert.NoError(t, err)
	l.Info("this is a test")
	if f, ok := (l.w.Writer).(*os.File); ok {
		assert.NoError(t, f.Close())
	} else {
		t.Error("Should be a file")
	}
	b, err := ioutil.ReadFile("test.log")
	assert.NoError(t, err)
	s := string(b)
	assert.Contains(t, s, "this is a test")
	assert.Contains(t, s, "INFO")
	assert.NoError(t, os.Remove("test.log"))
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
