package pkg

import (
	"fmt"
	"testing"
)

func TestUCFirst(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{name: "already uppercase", arg: "Number user count", want: "Number user count"},
		{name: "normal", arg: "userCount", want: "UserCount"},
	}

	for _, test := range tests {
		t.Run(test.name, func(subTest *testing.T) {
			if got := UCFirst(test.arg); got != test.want {
				subTest.Errorf("UCFirst, expect %v, got %v", test.want, got)
			}
		})
	}
}

func ExampleCamelCase() {
	fmt.Println(CamelCase("this_is_a_converter"))
	// Output: thisIsAConverter
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{name: "raw to camel", arg: "number user count", want: "numberUserCount"},
		{name: "snake to camel", arg: "number_user_count", want: "numberUserCount"},
		{name: "kebab to camel", arg: "number-user-count", want: "numberUserCount"},
		{name: "camel to camel", arg: "numberUserCount", want: "numberUserCount"},
		{name: "mix to camel", arg: "number-userCount", want: "numberUserCount"},
	}

	for _, test := range tests {
		t.Run(test.name, func(subTest *testing.T) {
			if got := CamelCase(test.arg); got != test.want {
				subTest.Errorf("CamelCase, expect %v, got %v", test.want, got)
			}
		})
	}
}

func ExampleSnakeCase() {
	fmt.Println(SnakeCase("ThisIsASnakeCaseConverter"))
	// Output: this_is_a_snake_case_converter
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{name: "raw to snake", arg: "number user count", want: "number_user_count"},
		{name: "snake to snake", arg: "ID0Value", want: "id0_value"},
		{name: "snake to snake", arg: "number_user_count", want: "number_user_count"},
		{name: "kebab to snake", arg: "number-user-count", want: "number_user_count"},
		{name: "camel to snake", arg: "numberUserCount", want: "number_user_count"},
		{name: "camel to snake", arg: "HTTPRequest", want: "http_request"},
		{name: "camel to snake", arg: "iN", want: "i_n"},
		{name: "camel to snake", arg: "INumber", want: "i_number"},
		{name: "camel to snake", arg: "Number", want: "number"},
		{name: "camel to snake", arg: "DNA", want: "dna"},
		{name: "camel to snake", arg: "numberDNASample", want: "number_dna_sample"},
		{name: "mix to snake", arg: "number-userCount", want: "number_user_count"},
	}

	for _, test := range tests {
		t.Run(test.name, func(subTest *testing.T) {
			if got := SnakeCase(test.arg); got != test.want {
				subTest.Errorf("SnakeCase, expect %v, got %v", test.want, got)
			}
		})
	}
}

func ExampleKebabCase() {
	fmt.Println(KebabCase("ThisIsAKebabCaseConverter"))
	// Output: this-is-a-kebab-case-converter
}

func TestKebabCase(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{name: "raw to kebab", arg: "number user count", want: "number-user-count"},
		{name: "snake to kebab", arg: "number_user_count", want: "number-user-count"},
		{name: "kebab to kebab", arg: "number-user-count", want: "number-user-count"},
		{name: "camel to kebab", arg: "numberUserCount", want: "number-user-count"},
		{name: "mix to kebab", arg: "number_userCount", want: "number-user-count"},
		{name: "mix to kebab 2", arg: "MigrateOldAccountHandler", want: "migrate-old-account-handler"},
	}

	for _, test := range tests {
		t.Run(test.name, func(subTest *testing.T) {
			if got := KebabCase(test.arg); got != test.want {
				subTest.Errorf("KebabCase, expect %v, got %v", test.want, got)
			}
		})
	}
}
