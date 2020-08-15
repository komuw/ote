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
		v := v // capture range variable
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
		add      bool
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
			add:      true,
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
			add:      true,
		},

		{
			line: &modfile.Line{
				Comments: modfile.Comments{
					// a `test` comment is added
					Suffix: []modfile.Comment{},
				},
			},
			expected: "// test",
			add:      true,
		},

		{
			line: &modfile.Line{
				Comments: modfile.Comments{
					Suffix: []modfile.Comment{
						{
							// existing `test` comment is REMOVED and any other comment left intact
							Token:  "// test; priorComment\n",
							Suffix: true,
						},
					},
				},
			},
			expected: "// priorComment\n",
			add:      false,
		},

		{
			line: &modfile.Line{
				Comments: modfile.Comments{
					Suffix: []modfile.Comment{
						{
							// existing `test` comment is removed
							Token:  "// test\n",
							Suffix: true,
						},
					},
				},
			},
			expected: "",
			add:      false,
		},
	}

	for _, v := range tt {
		v := v // capture range variable
		setTest(v.line, v.add)

		token := ""
		if v.line.Comments.Suffix != nil {
			// v.line.Comments.Suffix  is nil if there is no comment
			token = v.line.Comments.Suffix[0].Token
		}

		if !cmp.Equal(token, v.expected) {
			t.Errorf("\ngot \n\t%#+v \nwanted \n\t%#+v", token, v.expected)
		}
	}
}
