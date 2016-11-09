package seafile

import (
    "testing"
    "fmt"
    C "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { C.TestingT(t) }

type TestSuite struct{}

var _ = C.Suite(&TestSuite{})

const (
    userName = "choldrim@foxmail.com"
    password = "linuxdeepin123"
)

func (*TestSuite) TestGetToken(c *C.C) {
    token, err := GetToken(userName, password)
    c.Assert(err, C.Equals, nil)
    c.Assert(len(token), C.Not(C.Equals), 0)
}

func (*TestSuite) TestNewFolder(c *C.C) {
    s := New(userName, password)
    c.Assert(err, C.Equals, nil)
    c.Assert(len(token), C.Not(C.Equals), 0)
}
