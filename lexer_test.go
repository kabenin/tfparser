package tfparser

import (
	"fmt"
	"testing"
)

func TestSkipTillEOL(t *testing.T) {
	testStr := `# this is comment
This is first line after comment!
`

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	p.skipTillEOL()
	if p.data[p.i] != 'T' {
		t.Fatalf("After skiping comment next symbol is %#q, expected 'T'", p.data[p.i])
	}
}

func TestSkipMulitlineComment(t *testing.T) {
	testStr := `/*
	this is multiline comment
*/This is after`
	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	p.skipMulitlineComment()
	if p.data[p.i] != 'T' {
		t.Fatalf("After skiping multiline comment next symbol is %#q, expected 'T'", p.data[p.i])
	}
}

func TestSkipComment(t *testing.T) {
	testStr := `// Another
This is after`

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	p.skipComment()
	if p.data[p.i] != 'T' {
		t.Fatalf("After skiping comment next symbol is %#q, expected 'T'", p.data[p.i])
	}
}

func TestSkipCommentMultilineWithNewline(t *testing.T) {
	testStr := `/* mulitline */
T`

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	p.skipComment()
	if p.data[p.i] != 'T' {
		t.Fatalf("After skiping comment next symbol is %#q, expected 'T'", p.data[p.i])
	}
}

func TestSkipCommentMultilineWithNested(t *testing.T) {
	testStr := `/* mulitline
		/* nested **/
	*/
T`

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	p.skipComment()
	if p.data[p.i] != 'T' {
		t.Fatalf("After skiping comment next symbol is %#q, expected 'T'", p.data[p.i])
	}
}

func TestSkipCommentMultilineUnbalanced(t *testing.T) {
	testStr := `/* mulitline
		/* nested **/

T`

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	err := p.skipComment()
	if err == nil {
		t.Fatalf("skipComment did not return error on unbalanced comment")
	}
	if err.Error() != "Unable to find closing multiline comment" {
		t.Fatalf("Unexpected error returned by skipComment for unbalanced comment: %#q", err)
	}
}

func TestSkipBlockSimple(t *testing.T) {
	testStr := `{
		simple block
		}T`

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	p.skipBlock()
	if p.data[p.i] != 'T' {
		t.Fatalf("After skiping block next symbol is %#q, expected 'T'", p.data[p.i])
	}
}

func TestSkipBlockNested(t *testing.T) {
	testStr := `{
		another block {
			values
		}
	}T`

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	p.skipBlock()
	if p.data[p.i] != 'T' {
		t.Fatalf("After skiping block with nested next symbol is %#q, expected 'T'", p.data[p.i])
	}
}

func TestSkipBlockWithCommentedSubblock(t *testing.T) {
	testStr := `{
		/* another block {
			values
		}*/
	}T`

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	p.skipBlock()
	if p.data[p.i] != 'T' {
		t.Fatalf("After skiping block with commented out block next symbol is %#q, expected 'T'", p.data[p.i])
	}
}

func TestSkipBlockUnbalancedReturnsError(t *testing.T) {
	testStr := `{
		/* another block {
			values
		}*/
	T`

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	err := p.skipBlock()
	if err == nil {
		t.Fatalf("skipBlock did not return error with unbalanced block")
	}
	if err.Error() != "Unable to find closing brace for block" {
		t.Fatalf("skipBlock returned unexpected error for unbalanced block: %#q", err)
	}
}

func TestPeekQuotedStringWithLengthSimple(t *testing.T) {
	expected := "mytoken"
	testStr := fmt.Sprintf(`"%v" }
	`, expected)

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	token, l := p.peekQuotedStringWithLength()
	if p.err != nil {
		t.Fatalf("peekQuotedStringWithLength set an error: %q", p.err)
	}
	if token != expected {
		t.Fatalf("peekQuotedStringWithLength returned unexpected token %#q, expected %q", token, expected)
	}
	if l != len(expected)+2 {
		t.Fatalf("peekQuotedStringWithLength returned invalid length: %v, expected: %v", l, len(expected)+2)
	}
}

func TestPeekIdentifierWithLengthSimple(t *testing.T) {
	expected := "mytoken"
	testStr := fmt.Sprintf("%v =", expected)
	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	token, l := p.peekIdentifierWithLength()
	if p.err != nil {
		t.Fatalf("peekIdentifierWithLength triggered an error: %v", p.err)
	}
	if token != expected {
		t.Fatalf("peekIdentiferWithLength retunrned unexpected token %#q, expected %#q", token, expected)
	}
	if l != len(expected) {
		t.Fatalf("peekIdentifierWithLength returned unexpected length of a token: %v, expected %v", l, len(expected))
	}
}

func TestPeekSimpleKeyword(t *testing.T) {
	expected := "module"
	testStr := fmt.Sprintf(`
	%v "mytest`, expected)

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	token := p.peek()
	if token != expected {
		t.Fatalf("Unexpected token peeked: %#q, expected %#q", token, expected)
	}
}

func TestPeekSimpleQuotedIdentifier(t *testing.T) {
	expected := "test_module"
	testStr := fmt.Sprintf(`
	"%v" "mytest`, expected)

	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	token := p.peek()
	if token != expected {
		t.Fatalf("Unexpected token peeked: %#q, expected %#q", token, expected)
	}
}

func TestPeekSimpleQuotedIdentifierWithEscape(t *testing.T) {
	expected := "test_\\\"module"
	testStr := fmt.Sprint(`
	"test_\"module" "mytest`, expected)
	p := &parser{testStr, 0, &TFconfig{}, 0, nil, "", ""}
	token := p.peek()
	if token != expected {
		t.Fatalf("Unexpected token peeked: %#q, expected %#q", token, expected)
	}
}

var moduleTestData1 = `
	locals {
		variable = value
	}
	module "mytest_module" {
		var1 = "value one"
		var2 = "value two"

		providers = {
			aws1 = aws.us-west-1
			aws2 = aws.ap-southeast-2
		}
	}
`

func TestSomePops(t *testing.T) {
	p := &parser{moduleTestData1, 0, &TFconfig{}, 0, nil, "", ""}
	tokens := [...][2]string{
		{"locals", "{"},
		{"{", "v"},
		{"variable", "="},
		{"=", "v"},
		{"value", "}"},
		{"}", "m"},
		{"module", "\""},
		{"mytest_module", "{"},
		{"{", "v"},
		{"var1", "="},
		{"=", "\""},
		{"value one", "v"},
		{"var2", "="},
		{"=", "\""},
		{"value two", "p"},
		{"providers", "="},
		{"=", "{"},
		{"{", "a"},
		{"aws1", "="},
		{"=", "a"},
		{"aws.us-west-1", "a"},
		{"aws2", "="},
		{"=", "a"},
		{"aws.ap-southeast-2", "}"},
		{"}", "}"},
	}

	var probe [2]string
	for _, probe = range tokens {
		token := p.pop()

		if token != probe[0] {
			t.Fatalf("Unexpected token popped: %#q, expected %#q", token, probe[0])
		}
		if string(p.data[p.i]) != probe[1] {
			t.Fatalf("pop did not correctly advanced parsing position, epxected to be at %#q, but wer are at %#q", probe[1], p.data[p.i:len(p.data)-p.i])
		}
	}
	token := p.pop()
	if token != "}" {
		t.Fatalf("pop failed to pop final closing }. Got %#q", token)
	}
	if p.i < len(p.data) {
		t.Fatalf("after last pop in data index did not went out of boundaries: %v, %v", p.i, len(p.data))
	}
}

func TestSomePeeks(t *testing.T) {
	p := &parser{moduleTestData1, 0, &TFconfig{}, 0, nil, "", ""}
	tokens := [...][2]string{
		{"locals", "l"},
		{"{", "{"},
		{"variable", "v"},
		{"=", "="},
		{"value", "v"},
		{"}", "}"},
		{"module", "m"},
		{"mytest_module", "\""},
	}

	var probe [2]string
	for _, probe = range tokens {
		token := p.peek()
		if token != probe[0] {
			t.Fatalf("Unexpected token peeked: %#q, expected %#q", token, probe[0])
		}
		if string(p.data[p.i]) != probe[1] {
			t.Fatalf("peek unexpectedly advanced parsing position, epxected to be at %#q, but wer are at %#q", probe[1], p.data[p.i:len(p.data)-p.i])
		}
		p.pop()

	}
}
