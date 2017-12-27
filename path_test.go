package safepath

import (
	"testing"

	"encoding/json"

	. "gopkg.in/check.v1"
)

func TestPath(t *testing.T) {
	_ = Suite(&PathSuite{})
	TestingT(t)
}

type PathSuite struct {
}

func (s *PathSuite) TestUnsafeNewPath(c *C) {
	invalidButTrustedPath := "///../../../.."
	p := UnsafeNewPath(invalidButTrustedPath)
	c.Assert(p.raw, Equals, invalidButTrustedPath)
	c.Assert(p.String(), Equals, invalidButTrustedPath)
}

func (s *PathSuite) TestUnsafeNewRelativePath(c *C) {
	invalidButTrustedPath := "///../../../.."
	p := UnsafeNewRelativePath(invalidButTrustedPath)
	c.Assert(p.raw, Equals, invalidButTrustedPath)
	c.Assert(p.String(), Equals, invalidButTrustedPath)
}

func (s *PathSuite) TestSafeNewRelativePathBasic(c *C) {
	p, err := NewRelativePath("test")
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "test")
	c.Assert(p.String(), Equals, "test")
}

func (s *PathSuite) TestSafeNewRelativePathSubdir(c *C) {
	p, err := NewRelativePath("a/b/c")
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "a/b/c")
	c.Assert(p.String(), Equals, "a/b/c")
}

func (s *PathSuite) TestSafeNewRelativePathSubdirLongerNames(c *C) {
	p, err := NewRelativePath("abc/def/ghi")
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "abc/def/ghi")
	c.Assert(p.String(), Equals, "abc/def/ghi")
}

func (s *PathSuite) TestSafeNewRelativeDotted(c *C) {
	p, err := NewRelativePath(".a")
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, ".a")
	c.Assert(p.String(), Equals, ".a")
}

func (s *PathSuite) TestSafeNewRelativePathSubdirDotted(c *C) {
	p, err := NewRelativePath("a/.b")
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "a/.b")
	c.Assert(p.String(), Equals, "a/.b")
}

func (s *PathSuite) TestSafeNewRelativePathSubdirDotDotted(c *C) {
	p, err := NewRelativePath("a/..b")
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "a/..b")
	c.Assert(p.String(), Equals, "a/..b")
}

func (s *PathSuite) TestSafeNewRelativePathSubdirDotDotDotted(c *C) {
	p, err := NewRelativePath("a/...")
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "a/...")
	c.Assert(p.String(), Equals, "a/...")
}

func (s *PathSuite) TestSafeNewRelativePathSubdirDotDotDottedSubdir(c *C) {
	p, err := NewRelativePath("a/.../b")
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "a/.../b")
	c.Assert(p.String(), Equals, "a/.../b")
}

func (s *PathSuite) TestSafeNewRelativePathEmpty(c *C) {
	p, err := NewRelativePath("")
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "")
	c.Assert(p.String(), Equals, ".") // (!)
}

func (s *PathSuite) TestSafeNewRelativePathErrorAbsolute(c *C) {
	_, err := NewRelativePath("/")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathErrorAbsoluteFile(c *C) {
	_, err := NewRelativePath("/a")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathErrorAbsoluteDirectory(c *C) {
	_, err := NewRelativePath("/a/b")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathErrorSlashSlash(c *C) {
	_, err := NewRelativePath("a//b")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathErrorDot(c *C) {
	_, err := NewRelativePath(".")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathErrorDotDot(c *C) {
	_, err := NewRelativePath("..")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathErrorFileDot(c *C) {
	_, err := NewRelativePath("a/.")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathErrorFileDotDot(c *C) {
	_, err := NewRelativePath("a/..")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathErrorSubdirDot(c *C) {
	_, err := NewRelativePath("a/./b")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathErrorSubdirDotDot(c *C) {
	_, err := NewRelativePath("a/../b")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathErrorNullByte(c *C) {
	_, err := NewRelativePath("a/b\x00c/d")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestSafeNewRelativePathTrailingSlash(c *C) {
	_, err := NewRelativePath("a/b/")
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestRelativeJoinRelative(c *C) {
	l, err := NewRelativePath("a/b")
	c.Assert(err, IsNil)
	r, err := NewRelativePath("c/d")
	c.Assert(err, IsNil)
	p := l.Join(r)
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "a/b/c/d")
}

func (s *PathSuite) TestEmptyJoinRelative(c *C) {
	l, err := NewRelativePath("")
	c.Assert(err, IsNil)
	r, err := NewRelativePath("c/d")
	c.Assert(err, IsNil)
	p := l.Join(r)
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "c/d")
}

func (s *PathSuite) TestRelativeJoinEmpty(c *C) {
	l, err := NewRelativePath("a/b")
	c.Assert(err, IsNil)
	r, err := NewRelativePath("")
	c.Assert(err, IsNil)
	p := l.Join(r)
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "a/b")
}

func (s *PathSuite) TestEmptyRelativePathJoinEmptyRelativePath(c *C) {
	l, err := NewRelativePath("")
	c.Assert(err, IsNil)
	r, err := NewRelativePath("")
	c.Assert(err, IsNil)
	p := l.Join(r)
	c.Assert(err, IsNil)
	c.Assert(p.raw, Equals, "")
	c.Assert(p.String(), Equals, ".")
}

func (s *PathSuite) TestEmptyPathJoinEmptyRelativePath(c *C) {
	l := UnsafeNewPath("")
	r := UnsafeNewRelativePath("")
	p := l.Join(r)
	c.Assert(p.raw, Equals, "")
	c.Assert(p.String(), Equals, ".")
}

func (s *PathSuite) TestPathJoinRelativePath(c *C) {
	l := UnsafeNewPath("a/b")
	r := UnsafeNewRelativePath("c/d")
	p := l.Join(r)
	c.Assert(p.raw, Equals, "a/b/c/d")
	c.Assert(p.String(), Equals, "a/b/c/d")
}

func (s *PathSuite) TestJoinUnsafe(c *C) {
	l := UnsafeNewPath("a/b")
	p := l.JoinUnsafe("c/d")
	c.Assert(p.raw, Equals, "a/b/c/d")
	c.Assert(p.String(), Equals, "a/b/c/d")
}

func (s *PathSuite) TestEmptyJoinUnsafe(c *C) {
	l := UnsafeNewPath("")
	p := l.JoinUnsafe("c/d")
	c.Assert(p.raw, Equals, "c/d")
	c.Assert(p.String(), Equals, "c/d")
}

func (s *PathSuite) TestPathMarshalJSON(c *C) {
	l := UnsafeNewPath("/a/b/c")
	b, err := l.MarshalJSON()
	c.Assert(err, IsNil)
	c.Assert(b, DeepEquals, []byte("\"/a/b/c\""))
}

func (s *PathSuite) TestRelativePathMarshalJSON(c *C) {
	l := UnsafeNewRelativePath("a/b/c")
	b, err := l.MarshalJSON()
	c.Assert(err, IsNil)
	c.Assert(b, DeepEquals, []byte("\"a/b/c\""))
}

func (s *PathSuite) TestRelativePathUnmarshalJSON(c *C) {
	b := []byte("\"a/b/c\"")
	var l RelativePath
	err := json.Unmarshal(b, &l)
	c.Assert(err, IsNil)
	c.Assert(l.raw, Equals, "a/b/c")
}

func (s *PathSuite) TestRelativePathUnmarshalBadJSON(c *C) {
	b := []byte("{}")
	var l RelativePath
	err := json.Unmarshal(b, &l)
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestRelativePathUnmarshalJSONUnsafeValue(c *C) {
	b := []byte("\"/etc/passwd\"")
	var l RelativePath
	err := json.Unmarshal(b, &l)
	c.Assert(err, NotNil)
}

func (s *PathSuite) TestPathIsEmpty(c *C) {
	p := UnsafeNewPath("")
	c.Assert(p.IsEmpty(), Equals, true)
}

func (s *PathSuite) TestPathIsNotEmpty(c *C) {
	p := UnsafeNewPath("aaa")
	c.Assert(p.IsEmpty(), Equals, false)
}

func (s *PathSuite) TestRelativePathIsEmpty(c *C) {
	p := UnsafeNewRelativePath("")
	c.Assert(p.IsEmpty(), Equals, true)
}

func (s *PathSuite) TestRelativePathIsNotEmpty(c *C) {
	p := UnsafeNewRelativePath("aaa")
	c.Assert(p.IsEmpty(), Equals, false)
}

func (s *PathSuite) TestFilePathBase(c *C) {
	p := UnsafeNewRelativePath("aa")
	c.Assert(p.Base(), Equals, "aa")
}

func (s *PathSuite) TestSubdirFilePathBase(c *C) {
	p := UnsafeNewRelativePath("aa/bb/cc")
	c.Assert(p.Base(), Equals, "cc")
}

func (s *PathSuite) TestEmptyPathBase(c *C) {
	p := UnsafeNewRelativePath("")
	c.Assert(p.Base(), Equals, "")
}
