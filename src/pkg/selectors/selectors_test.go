package selectors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SelectorSuite struct {
	suite.Suite
}

func TestSelectorSuite(t *testing.T) {
	suite.Run(t, new(SelectorSuite))
}

func (suite *SelectorSuite) TestNewSelector() {
	t := suite.T()
	s := newSelector(ServiceUnknown)
	assert.NotNil(t, s)
	assert.Equal(t, s.Service, ServiceUnknown)
	assert.NotNil(t, s.Includes)
}

func (suite *SelectorSuite) TestBadCastErr() {
	err := badCastErr(ServiceUnknown, ServiceExchange)
	assert.Error(suite.T(), err)
}

func (suite *SelectorSuite) TestPrintable() {
	t := suite.T()

	sel := stubSelector()
	p := sel.Printable()

	assert.Equal(t, sel.Service.String(), p.Service)
	assert.Equal(t, 1, len(p.Excludes))
	assert.Equal(t, 1, len(p.Filters))
	assert.Equal(t, 1, len(p.Includes))
}

func (suite *SelectorSuite) TestPrintable_IncludedResources() {
	t := suite.T()

	sel := stubSelector()
	p := sel.Printable()
	res := p.Resources()

	assert.Equal(t, "All", res, "stub starts out as an all-pass")

	stubWithResource := func(resource string) scope {
		ss := stubScope("")
		ss[rootCatStub.String()] = filterize(scopeConfig{}, resource)

		return scope(ss)
	}

	sel.Includes = []scope{
		stubWithResource("foo"),
		stubWithResource("smarf"),
		stubWithResource("fnords"),
	}

	p = sel.Printable()
	res = p.Resources()

	assert.True(t, strings.HasSuffix(res, "(2 more)"), "resource '"+res+"' should have (2 more) suffix")

	p.Includes = nil
	res = p.Resources()

	assert.Equal(t, "All", res, "filters is also an all-pass")

	p.Filters = nil
	res = p.Resources()

	assert.Equal(t, "None", res, "resource with no Includes or Filters should state None")
}

func (suite *SelectorSuite) TestToResourceTypeMap() {
	table := []struct {
		name   string
		input  []scope
		expect map[string][]string
	}{
		{
			name:  "single scope",
			input: []scope{scope(stubScope(""))},
			expect: map[string][]string{
				"All": {rootCatStub.String()},
			},
		},
		{
			name: "disjoint resources",
			input: []scope{
				scope(stubScope("")),
				{
					rootCatStub.String(): filterize(scopeConfig{}, "smarf"),
					scopeKeyDataType:     filterize(scopeConfig{}, unknownCatStub.String()),
				},
			},
			expect: map[string][]string{
				"All":   {rootCatStub.String()},
				"smarf": {unknownCatStub.String()},
			},
		},
		{
			name: "disjoint types",
			input: []scope{
				scope(stubScope("")),
				{
					rootCatStub.String(): filterize(scopeConfig{}, AnyTgt),
					scopeKeyDataType:     filterize(scopeConfig{}, "other"),
				},
			},
			expect: map[string][]string{
				"All": {rootCatStub.String(), "other"},
			},
		},
	}
	for _, test := range table {
		suite.T().Run(test.name, func(t *testing.T) {
			rtm := toResourceTypeMap[mockScope](test.input)
			assert.Equal(t, test.expect, rtm)
		})
	}
}

func (suite *SelectorSuite) TestContains() {
	t := suite.T()
	key := rootCatStub
	target := "fnords"
	does := stubScope("")
	does[key.String()] = filterize(scopeConfig{}, target)
	doesNot := stubScope("")
	doesNot[key.String()] = filterize(scopeConfig{}, "smarf")

	assert.True(t, matches(does, key, target), "does contain")
	assert.False(t, matches(doesNot, key, target), "does not contain")
}
