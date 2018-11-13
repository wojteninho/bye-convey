package bye_convey

import (
	"bytes"
)

type transformFn func(in []byte) []byte

var transformations = []transformFn{removeImportFn, replaceConveyWithTestingRunFn, replaceTestingRunFn, multiplyClosingParenthesisFn, replaceAssertionsFn}

func Transform(in []byte) []byte {
	for _, t := range transformations {
		in = t(in)
	}

	return in
}

func removeImportFn(in []byte) []byte {
	if !bytes.Equal(in, []byte("	. \"github.com/smartystreets/goconvey/convey\"")) {
		return in
	}

	return bytes.Replace(in, []byte("	. \"github.com/smartystreets/goconvey/convey\""), []byte("	\"github.com/Propertyfinder/pf-go-common/pkg/test\"\n	. \"github.com/onsi/gomega\""), 1)
}

func replaceConveyWithTestingRunFn(in []byte) []byte {
	if !bytes.Contains(in, []byte("Convey")) {
		return in
	}

	// Convey --> t.Run
	in = bytes.Replace(in, []byte("Convey("), []byte("t.Run("), 1)

	// t, func() { --> func(t *testing.T) {
	in = bytes.Replace(in, []byte("t, func() {"), []byte("test.GomegaWrapper(func(t *testing.T) {"), 1)

	return in
}

func replaceTestingRunFn(in []byte) []byte {
	if !bytes.Contains(in, []byte("t.Run")) {
		return in
	}

	if !bytes.Contains(in, []byte(", func() {")) {
		return in
	}

	// func() { --> func(t *testing.T) {
	return bytes.Replace(in, []byte(", func() {"), []byte(", test.GomegaWrapper(func(t *testing.T) {"), 1)
}

/**
When we replace:

t, func() { -> test.GomegaWrapper(func(t *testing.T) {
, func() { -> , test.GomegaWrapper(func(t *testing.T) {

we miss one closing parenthesis. W need to detect the line which contains only }) and replace it with })).
 */
func multiplyClosingParenthesisFn(in []byte) []byte {
	if !bytes.Equal(bytes.Replace(bytes.TrimSpace(in),[]byte("\t"), []byte(""),-1), []byte("})")) {
		return in
	}

	// }) --> }))
	return bytes.Replace(in, []byte("})"), []byte("}))"), 1)
}

func replaceAssertionsFn(in []byte) []byte {
	// So --> Expect
	in = bytes.Replace(in, []byte("So("), []byte("Expect("), 1)

	// ---------------------------------------------------
	// [x] ShouldEqual          = assertions.ShouldEqual
	// [x] ShouldNotEqual       = assertions.ShouldNotEqual
	// [ ] ShouldAlmostEqual    = assertions.ShouldAlmostEqual
	// [ ] ShouldNotAlmostEqual = assertions.ShouldNotAlmostEqual
	// [x] ShouldResemble       = assertions.ShouldResemble
	// [x] ShouldNotResemble    = assertions.ShouldNotResemble
	// [ ] ShouldPointTo        = assertions.ShouldPointTo
	// [ ] ShouldNotPointTo     = assertions.ShouldNotPointTo
	// [x] ShouldBeNil          = assertions.ShouldBeNil
	// [x] ShouldNotBeNil       = assertions.ShouldNotBeNil
	// [x] ShouldBeTrue         = assertions.ShouldBeTrue
	// [x] ShouldBeFalse        = assertions.ShouldBeFalse
	// [ ] ShouldBeZeroValue    = assertions.ShouldBeZeroValue

	// ShouldEqual --> To(Equal(...))
	in = bytes.Replace(in, []byte(", ShouldEqual, "), []byte(").To(Equal("), 1)
	in = bytes.Replace(in, []byte(", ShouldNotEqual, "), []byte(").ToNot(Equal("), 1)

	// ShouldResemble -->
	in = bytes.Replace(in, []byte(", ShouldResemble, "), []byte(").To(Equal("), 1)
	in = bytes.Replace(in, []byte(", ShouldNotResemble, "), []byte(").ToNot(Equal("), 1)

	// ShouldBeNil --> To(BeNil())
	in = bytes.Replace(in, []byte(", ShouldBeNil)"), []byte(").To(BeNil())"), 1)
	in = bytes.Replace(in, []byte(", ShouldNotBeNil)"), []byte(").ToNot(BeNil())"), 1)

	// ShouldBeTrue --> To(BeTrue())
	in = bytes.Replace(in, []byte(", ShouldBeTrue)"), []byte(").To(BeTrue())"), 1)

	// ShouldBeFalse --> To(BeFalse())
	in = bytes.Replace(in, []byte(", ShouldBeFalse)"), []byte(").To(BeFalse())"), 1)

	// -------------------------------------------------------
	// [ ] ShouldBeGreaterThan          = assertions.ShouldBeGreaterThan
	// [ ] ShouldBeGreaterThanOrEqualTo = assertions.ShouldBeGreaterThanOrEqualTo
	// [ ] ShouldBeLessThan             = assertions.ShouldBeLessThan
	// [ ] ShouldBeLessThanOrEqualTo    = assertions.ShouldBeLessThanOrEqualTo
	// [ ] ShouldBeBetween              = assertions.ShouldBeBetween
	// [ ] ShouldNotBeBetween           = assertions.ShouldNotBeBetween
	// [ ] ShouldBeBetweenOrEqual       = assertions.ShouldBeBetweenOrEqual
	// [ ] ShouldNotBeBetweenOrEqual    = assertions.ShouldNotBeBetweenOrEqual

	// -------------------------------------------------------
	// [x] ShouldContain       = assertions.ShouldContain
	// [x] ShouldNotContain    = assertions.ShouldNotContain
	// [x] ShouldContainKey    = assertions.ShouldContainKey
	// [x] ShouldNotContainKey = assertions.ShouldNotContainKey
	// [ ] ShouldBeIn          = assertions.ShouldBeIn
	// [ ] ShouldNotBeIn       = assertions.ShouldNotBeIn
	// [x] ShouldBeEmpty       = assertions.ShouldBeEmpty
	// [x] ShouldNotBeEmpty    = assertions.ShouldNotBeEmpty
	// [x] ShouldHaveLength    = assertions.ShouldHaveLength

	// ShouldContainKey --> To(HaveKey(...))
	in = bytes.Replace(in, []byte(", ShouldContain, "), []byte(").To(ContainElement("), 1)
	in = bytes.Replace(in, []byte(", ShouldNotContain, "), []byte(").ToNot(ContainElement("), 1)

	// ShouldContainKey --> To(HaveKey(...))
	in = bytes.Replace(in, []byte(", ShouldContainKey, "), []byte(").To(HaveKey("), 1)
	in = bytes.Replace(in, []byte(", ShouldNotContainKey, "), []byte(").ToNot(HaveKey("), 1)

	// ShouldBeEmpty  --> To(BeEmpty())
	in = bytes.Replace(in, []byte(", ShouldBeEmpty)"), []byte(").To(BeEmpty())"), 1)
	in = bytes.Replace(in, []byte(", ShouldNotBeEmpty)"), []byte(").ToNot(BeEmpty())"), 1)

	// ShouldHaveLength  --> To(HaveLen(...))
	in = bytes.Replace(in, []byte(", ShouldHaveLength, "), []byte(").To(HaveLen("), 1)

	// -------------------------------------------------------
	// [ ] ShouldStartWith           = assertions.ShouldStartWith
	// [ ] ShouldNotStartWith        = assertions.ShouldNotStartWith
	// [ ] ShouldEndWith             = assertions.ShouldEndWith
	// [ ] ShouldNotEndWith          = assertions.ShouldNotEndWith
	// [ ] ShouldBeBlank             = assertions.ShouldBeBlank
	// [ ] ShouldNotBeBlank          = assertions.ShouldNotBeBlank
	// [x] ShouldContainSubstring    = assertions.ShouldContainSubstring
	// [x] ShouldNotContainSubstring = assertions.ShouldNotContainSubstring

	// ShouldContainSubstring  --> To(ContainSubstring(...))
	in = bytes.Replace(in, []byte(", ShouldContainSubstring, "), []byte(").To(ContainSubstring("), 1)
	in = bytes.Replace(in, []byte(", ShouldNotContainSubstring, "), []byte(").ToNot(ContainSubstring("), 1)

	// -------------------------------------------------------
	// [x] ShouldPanic        = assertions.ShouldPanic
	// [x] ShouldNotPanic     = assertions.ShouldNotPanic
	// [ ] ShouldPanicWith    = assertions.ShouldPanicWith
	// [ ] ShouldNotPanicWith = assertions.ShouldNotPanicWith

	// ShouldPanic --> To(Panic())
	in = bytes.Replace(in, []byte(", ShouldPanic)"), []byte(").To(Panic())"), 1)
	in = bytes.Replace(in, []byte(", ShouldNotPanic)"), []byte(").ToNot(Panic())"), 1)

	// -------------------------------------------------------
	// [x] ShouldHaveSameTypeAs    = assertions.ShouldHaveSameTypeAs
	// [x] ShouldNotHaveSameTypeAs = assertions.ShouldNotHaveSameTypeAs
	// [x] ShouldImplement         = assertions.ShouldImplement
	// [x] ShouldNotImplement      = assertions.ShouldNotImplement

	// ShouldHaveSameTypeAs --> To(Equal(...))
	in = bytes.Replace(in, []byte(", ShouldHaveSameTypeAs, "), []byte(").To(BeAssignableToTypeOf("), 1)
	in = bytes.Replace(in, []byte(", ShouldNotHaveSameTypeAs, "), []byte(").ToNot(BeAssignableToTypeOf("), 1)
	in = bytes.Replace(in, []byte(", ShouldImplement, "), []byte(").To(test.Implements("), 1)
	in = bytes.Replace(in, []byte(", ShouldNotImplement, "), []byte(").ToNot(test.Implements("), 1)

	// -------------------------------------------------------
	// [x] ShouldHappenBefore         = assertions.ShouldHappenBefore
	// [x] ShouldHappenOnOrBefore     = assertions.ShouldHappenOnOrBefore
	// [x] ShouldHappenAfter          = assertions.ShouldHappenAfter
	// [x] ShouldHappenOnOrAfter      = assertions.ShouldHappenOnOrAfter
	// [ ] ShouldHappenBetween        = assertions.ShouldHappenBetween
	// [ ] ShouldHappenOnOrBetween    = assertions.ShouldHappenOnOrBetween
	// [ ] ShouldNotHappenOnOrBetween = assertions.ShouldNotHappenOnOrBetween
	// [ ] ShouldHappenWithin         = assertions.ShouldHappenWithin
	// [ ] ShouldNotHappenWithin      = assertions.ShouldNotHappenWithin
	// [ ] ShouldBeChronological      = assertions.ShouldBeChronological

	// ShouldHappenBefore --> To(BeTemporally("<"))
	in = bytes.Replace(in, []byte(", ShouldHappenBefore, "), []byte(").To(BeTemporally(\"<\", "), 1)

	// ShouldHappenOnOrBefore --> To(BeTemporally("<="))
	in = bytes.Replace(in, []byte(", ShouldHappenOnOrBefore, "), []byte(").To(BeTemporally(\"<=\", "), 1)

	// ShouldHappenAfter --> To(BeTemporally(">"))
	in = bytes.Replace(in, []byte(", ShouldHappenAfter, "), []byte(").To(BeTemporally(\">\", "), 1)

	// ShouldHappenOnOrAfter --> To(BeTemporally(">="))
	in = bytes.Replace(in, []byte(", ShouldHappenOnOrAfter, "), []byte(").To(BeTemporally(\">=\", "), 1)

	// -------------------------------------------------------
	// [ ] ShouldBeError = assertions.ShouldBeError

	// ShouldBeError --> To(HaveOccurred())
	in = bytes.Replace(in, []byte(", ShouldBeError)"), []byte(").To(HaveOccurred())"), 1)

	// append one more )
	if  bytes.Contains(in, []byte("Equal(")) ||
		bytes.Contains(in, []byte("BeAssignableToTypeOf(")) ||
		bytes.Contains(in, []byte("Implements(")) ||
		bytes.Contains(in, []byte("HaveLen(")) ||
		bytes.Contains(in, []byte("HaveKey(")) ||
		bytes.Contains(in, []byte("BeTemporally(")) ||
		bytes.Contains(in, []byte("ContainElement(")) ||
		bytes.Contains(in, []byte("ContainSubstring(")) {
		in = append(in, []byte(")")...)
	}

	return in
}