package bridge

import (
	// "fmt"
	// "go/parser"
	// "go/token"
	. "gopkg.in/check.v1"
	// "log"
	// "testing"
)

// Tests begin

func (s *TestSuite) TestTypeClassString(c *C) {
	c.Assert(NewType(NullType, nil).TypeClassString(), Equals, "NullType")
	c.Assert(NewType(UnresolvedType, nil).TypeClassString(), Equals, "UnresolvedType")
	c.Assert(NewType(NamedType, nil).TypeClassString(), Equals, "NamedType")
	c.Assert(NewType(AliasType, nil).TypeClassString(), Equals, "AliasType")
	c.Assert(NewType(ReferenceType, nil).TypeClassString(), Equals, "ReferenceType")
	c.Assert(NewType(TupleType, nil).TypeClassString(), Equals, "TupleType")
	c.Assert(NewType(RecordType, nil).TypeClassString(), Equals, "RecordType")
	c.Assert(NewType(FunctionType, nil).TypeClassString(), Equals, "FunctionType")
	c.Assert(NewType(ListType, nil).TypeClassString(), Equals, "ListType")
	c.Assert(NewType(MapType, nil).TypeClassString(), Equals, "MapType")
}

func (s *TestSuite) TestIsValueType(c *C) {
	c.Assert(NewType(ReferenceType, nil).IsValueType(), Equals, false)
	c.Assert(NewType(FunctionType, nil).IsValueType(), Equals, false)
	c.Assert(NewType(ListType, nil).IsValueType(), Equals, false)
	c.Assert(NewType(MapType, nil).IsValueType(), Equals, false)

	c.Assert(NewType(NamedType, nil).IsValueType(), Equals, true)
	c.Assert(NewType(RecordType, nil).IsValueType(), Equals, true)
}

func (s *TestSuite) TestNewType(c *C) {
	cls := 10
	d := "Hello"
	t := NewType(cls, d)
	c.Assert(t, Not(IsNil))
	c.Assert(cls, Equals, t.TypeClass)
	c.Assert(d, Equals, t.TypeData)
}
