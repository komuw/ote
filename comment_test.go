package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/mod/modfile"
)

func TestIsTest(t *testing.T) {

	tt := []struct {
		line     *modfile.Line
		expected bool
	}{

		{
			line: &modfile.Line{
				Comments: modfile.Comments{
					Suffix: []modfile.Comment{
						{
							Token:  "// test\n",
							Suffix: true,
						},
					},
				},
				Token: []string{
					"github.com/stretchr/testify",
					"v1.6.1",
				},
			},
			expected: true,
		},

		{
			line: &modfile.Line{
				Comments: modfile.Comments{
					Suffix: []modfile.Comment{
						{
							Token:  "// cool\n",
							Suffix: true,
						},
					},
				},
				Token: []string{
					"github.com/pkg/errors",
					"v0.08",
				},
			},
			expected: false,
		},
	}

	for _, v := range tt {
		res := isTest(v.line)

		if !cmp.Equal(res, v.expected) {
			t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", v, v.expected)
		}

	}

}

func TestSetTest(t *testing.T) {

	tt := []struct {
		line     *modfile.Line
		expected string
	}{

		{
			line: &modfile.Line{
				Comments: modfile.Comments{
					Suffix: []modfile.Comment{
						{
							// existing `test` comment left intact
							Token:  "// test\n",
							Suffix: true,
						},
					},
				},
			},
			expected: "// test\n",
		},

		{
			line: &modfile.Line{
				Comments: modfile.Comments{
					Suffix: []modfile.Comment{
						{
							// existing comment left intact & `test` comment is added
							Token:  "// someOtherComment",
							Suffix: true,
						},
					},
				},
			},
			expected: "// test; someOtherComment",
		},

		{
			line: &modfile.Line{
				Comments: modfile.Comments{
					// a `test` comment is added
					Suffix: []modfile.Comment{},
				},
			},
			expected: "// test",
		},
	}

	for _, v := range tt {
		setTest(v.line)
		token := v.line.Comments.Suffix[0].Token

		if !cmp.Equal(token, v.expected) {
			t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", token, v.expected)
		}

	}

}
