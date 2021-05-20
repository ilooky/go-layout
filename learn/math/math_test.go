package math

import (
	"fmt"
	"testing"
)

//go test
//go test -v

//Writing Coverage Tests in Go
//go test -coverprofile=coverage.out
//go tool cover -html=coverage.out

func TestAdd(t *testing.T) {

	got := Add(4, 6)
	want := 10

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

// arg1 means argument 1 and arg2 means argument 2, and the expected stands for the 'result we expect'
type addTest struct {
	arg1, arg2, expected int
}

var addTests = []addTest{
	addTest{2, 3, 5},
	addTest{4, 8, 12},
	addTest{6, 9, 15},
	addTest{3, 10, 13},
}

func TestAddTable(t *testing.T) {
	for _, test := range addTests {
		if output := Add(test.arg1, test.arg2); output != test.expected {
			t.Errorf("Output %d not equal to expected %d", output, test.expected)
		}
	}
}

//go test -bench=.
//go test -bench=Add
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(4, 6)
	}
}

//go test -v
func ExampleAdd() {
	fmt.Println(Add(4, 6))
	// Output: 11
}