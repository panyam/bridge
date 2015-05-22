package bridge

import (
	// "fmt"
	// "go/parser"
	// "go/token"
	. "gopkg.in/check.v1"
	"log"
	"testing"
)

type TestSuite struct {
}

var _ = Suite(&TestSuite{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	log.Println("Setting up tests....")
	TestingT(t)
}

func (s *TestSuite) SetUpSuite(c *C) {
}

func (s *TestSuite) SetUpTest(c *C) {
}

func (s *TestSuite) TearDownTest(c *C) {
}

// Tests begin

func (s *TestSuite) TestAddNamedType(c *C) {
	ts := NewTypeLibrary()
	t := NewType(NamedType, nil)
	t = ts.AddType("", "int64", t)
	t2 := ts.GetType("", "int64")
	c.Assert(t2, Not(IsNil))
	c.Assert(t.TypeClass, Equals, t2.TypeClass)
	c.Assert(t.TypeData, Equals, t2.TypeData)
}
