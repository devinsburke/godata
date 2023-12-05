package godata

import (
	"context"
	"strings"
	"testing"
)

func TestParseCompute(t *testing.T) {
	err := DefineCustomFunctions([]CustomFunctionInput{
		{Name: "zeroArgFunc", NumParams: []int{0}},
		{Name: "oneArgFunc", NumParams: []int{1}},
		{Name: "twoArgFunc", NumParams: []int{2}},
	})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	type testcase struct {
		computeStrings []string
		shouldPass     bool
	}

	testCases := []testcase{
		{[]string{"oldField as newField"}, true},
		{[]string{"1 as newField"}, true},
		{[]string{"one add 2 as newField"}, true},
		{[]string{"one add two as extra/newField"}, true},
		{[]string{"zeroArgFunc() as newField"}, true},
		{[]string{"oneArgFunc(one) as newField"}, true},
		{[]string{"twoArgFunc(one, two) as newField"}, true},
		{[]string{"twoArgFunc(one, two) as newField", "tolower(three) as  newFieldTwo"}, true},

		{[]string{"case(false:0) as newField"}, true},
		{[]string{"case(false:0,true:1) as newField"}, true},
		{[]string{"case(prop eq 'one':1,prop eq 'two':2) as newField"}, true},
		{[]string{"case(tolower(one) eq one:'lower') as newField"}, true},
		{[]string{"case(contains(haystack,'needle'):1,true:1) as newField"}, true},
		{[]string{"case(tolower(one) eq one:tolower(one)) as newField"}, true},
		{[]string{"case(true:2 mul 3) as newField"}, true},
		{[]string{"case(false:1,false:2,false:3,false:4,false:5,false:6,false:7,false:8,false:9,false:10) as newField"}, true}, // max of 10 cases

		// negative cases
		{[]string{"one add two as newField2"}, false},
		{[]string{"one add two newField2"}, false},
		{[]string{""}, false},
		{[]string{"as"}, false},
		{[]string{"as newField"}, false},
		{[]string{"zeroArgFunc() as "}, false},

		{[]string{"case as newField"}, false},
		{[]string{"case() as newField"}, false},
		{[]string{"case(false:,true:1) as newField"}, false},
		{[]string{"case(false,true:1) as newField"}, false},
		{[]string{"case(false,true:1) as newField"}, false},
		{[]string{"case(1:1,true:1) as newField"}, false},
		{[]string{"case(1:1,true:1) as newField"}, false},
		{[]string{"case(false:1,false:2,false:3,false:4,false:5,false:6,false:7,false:8,false:9,false:10,false:11) as newField"}, false}, // max of 10 cases
	}

	for i, v := range testCases {

		var result *GoDataComputeQuery
		result, err = ParseComputeString(context.Background(), strings.Join(v.computeStrings, ","))
		if v.shouldPass {
			if err != nil {
				t.Errorf("testcase %d: %v", i, err)
				t.FailNow()
			}
			if result == nil {
				t.Errorf("testcase %d: nil result", i)
				t.FailNow()
			}

			if len(result.ComputeItems) != len(v.computeStrings) {
				t.Errorf("testcase %d: wrong number of result items", i)
				t.FailNow()
			}
		} else {
			if err == nil {
				t.Errorf("testcase %d: expected error", i)
				t.FailNow()
			}
		}

	}
}
