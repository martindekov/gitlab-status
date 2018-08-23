package function

import (
	"testing"
)

func test_getURLbyDelimeter(t *testing.T) {
	stringStruct := []struct {
		Title          string
		testString     string
		delimeter      string
		position       int
		expectedString string
	}{
		{
			"Delimeter is . and everything before the 4th place of the delimeter",
			"some.random.talking.here.",
			".",
			4,
			"some.random.talking",
		},
		{
			"This is base URL with delimeter / and everything before 3rd position",
			"https://some.random.url/some/routing/here",
			"/",
			3,
			"https://some.random.url",
		},
		{
			"When using letter z as delimeter everything before 5th place",
			"zzzz,I love sleeping,zzzz",
			"z",
			5,
			"zzzz,I love sleeping,",
		},
	}
	for _, test := range stringStruct {
		t.Run(test.Title, func(t *testing.T) {
			leftoverString := getURLbyDelimeter(test.testString, test.position, test.delimeter)
			if leftoverString != test.expectedString {
				t.Errorf("Expected string: `%s` got: `%s`", test.expectedString, leftoverString)
			}
		})
	}
}
