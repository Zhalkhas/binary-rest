package repos

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"slices"
	"strings"
	"testing"
)

func TestNewReaderIndices(t *testing.T) {
	tt := map[string]struct {
		input       io.Reader
		expected    []int64
		expectedErr error
	}{
		"should return empty array to empty input": {
			input:       strings.NewReader(""),
			expected:    []int64{},
			expectedErr: nil,
		},
		"should return error to invalid input": {
			input:       strings.NewReader("asjkdnaksd"),
			expected:    []int64{},
			expectedErr: ErrInvalidIndicesFile,
		},
		"should parse valid input": {
			input:       strings.NewReader("1\n2\n3\n4\n5"),
			expected:    []int64{1, 2, 3, 4, 5},
			expectedErr: nil,
		},
		"parsing error should be returned": {
			input:       mockFailingReader{},
			expected:    []int64{},
			expectedErr: ErrParsingIndicesFile,
		},
	}

	for name, tCase := range tt {
		t.Run(name, func(t *testing.T) {
			actual, actualErr := NewReaderIndices(tCase.input)
			if !errors.Is(actualErr, tCase.expectedErr) {
				t.Errorf("expected error %v, got %v", tCase.expectedErr, actualErr)
			}
			if !slices.Equal(actual.values, tCase.expected) {
				t.Errorf("expected %v, got %v", tCase.expected, actual.values)
			}
		})
	}
}

type mockFailingReader struct{}

func (m mockFailingReader) Read(_ []byte) (n int, err error) {
	return -1, nil
}

func TestReaderIndices_Search(t *testing.T) {
	tt := map[string]struct {
		inputArr       io.Reader
		inputSearchVal int64
		expected       int
		expectedErr    error
	}{
		"should return error to empty array": {
			inputArr:       strings.NewReader(""),
			inputSearchVal: 1,
			expected:       -1,
			expectedErr:    ErrIndexNotFound,
		},
		"should return error to not found value (arr length is odd)": {
			inputArr:       strings.NewReader("1\n2\n3\n4\n5"),
			inputSearchVal: 1000,
			expected:       -1,
			expectedErr:    ErrIndexNotFound,
		},
		"should return error to not found value (arr length is even)": {
			inputArr:       strings.NewReader("1\n2\n3\n4\n5\n6"),
			inputSearchVal: 1000,
			expected:       -1,
			expectedErr:    ErrIndexNotFound,
		},
		"should return index to found value (arr length odd)": {
			inputArr:       strings.NewReader("1\n2\n3\n4\n5"),
			inputSearchVal: 3,
			expected:       2,
			expectedErr:    nil,
		},
		"should return index to found value (arr length even)": {
			inputArr:       strings.NewReader("1\n2\n3\n4\n5\n6"),
			inputSearchVal: 3,
			expected:       2,
			expectedErr:    nil,
		},
		"should return index to found value (arr length odd) when searched value is first": {
			inputArr:       strings.NewReader("1\n2\n3\n4\n5"),
			inputSearchVal: 1,
			expected:       0,
			expectedErr:    nil,
		},
		"should return index to found value (arr length even) when searched value is first": {
			inputArr:       strings.NewReader("1\n2\n3\n4\n5\n6"),
			inputSearchVal: 1,
			expected:       0,
			expectedErr:    nil,
		}, "should return index to found value (arr length odd) when searched value is last": {
			inputArr:       strings.NewReader("1\n2\n3\n4\n5"),
			inputSearchVal: 5,
			expected:       4,
			expectedErr:    nil,
		},
		"should return index to found value (arr length even) when searched value is last": {
			inputArr:       strings.NewReader("1\n2\n3\n4\n5\n6"),
			inputSearchVal: 6,
			expected:       5,
			expectedErr:    nil,
		},
		"should return closest index to found value (arr length even)": {
			inputArr:       strings.NewReader("1000\n1100\n1200\n1300"),
			inputSearchVal: 1150,
			expected:       1,
			expectedErr:    nil,
		},
		"should return closest index to found value  (arr length odd)": {
			inputArr:       strings.NewReader("1000\n1100\n1200\n1300\n1400"),
			inputSearchVal: 1150,
			expected:       1,
			expectedErr:    nil,
		},
		"should return closest index to found value (arr length even) when closest is first": {
			inputArr:       strings.NewReader("1000\n1100\n1200\n1300"),
			inputSearchVal: 1050,
			expected:       0,
			expectedErr:    nil,
		}, "should return closest index to found value  (arr length odd) when closest is first": {
			inputArr:       strings.NewReader("1000\n1100\n1200\n1300\n1400"),
			inputSearchVal: 1050,
			expected:       0,
			expectedErr:    nil,
		},
		"should return closest index to found value (arr length even) when closest is last": {
			inputArr:       strings.NewReader("1000\n1100\n1200\n1300"),
			inputSearchVal: 1350,
			expected:       3,
			expectedErr:    nil,
		},
		"should return closest index to found value  (arr length odd) when closest is last": {
			inputArr:       strings.NewReader("1000\n1100\n1200\n1300\n1400"),
			inputSearchVal: 1450,
			expected:       4,
			expectedErr:    nil,
		}, "should return closest index to found value  (arr length odd) when closest is last 2": {
			inputArr:       strings.NewReader("1000\n1200\n1400"),
			inputSearchVal: 1350,
			expected:       2,
			expectedErr:    nil,
		},
		"should return error if value is not close enough": {
			inputArr:       strings.NewReader("1000\n1500\n2000"),
			inputSearchVal: 1150,
			expected:       -1,
			expectedErr:    ErrIndexNotFound,
		},
	}

	for name, tCase := range tt {
		t.Run(name, func(t *testing.T) {
			reader, err := NewReaderIndices(tCase.inputArr)
			if err != nil {
				t.Errorf("unexpected error creating reader: %v", err)
			}
			actual, actualErr := reader.Search(context.Background(), tCase.inputSearchVal)
			if !errors.Is(actualErr, tCase.expectedErr) {
				t.Errorf("expected error %v, got %v", tCase.expectedErr, actualErr)
			}
			if actual != tCase.expected {
				t.Errorf("expected %v, got %v", tCase.expected, actual)
			}
		})
	}
	t.Run("canceling context must return error", func(t *testing.T) {
		reader, err := NewReaderIndices(strings.NewReader(generateSortedInts(1000000)))
		if err != nil {
			t.Errorf("unexpected error creating reader: %v", err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			_, err := reader.Search(ctx, 1)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("expected error %v, got %v", context.Canceled, err)
			}
		}()
		cancel()
	})
}

func generateSortedInts(length int) string {
	input := make([]int, length)
	for i := range input {
		input[i] = i
	}
	inputStrs := make([]string, length)
	for i := range input {
		inputStrs[i] = fmt.Sprint(input[i])
	}
	return strings.Join(inputStrs, "\n")
}

func generateRandomSortedInts(length int) string {
	input := make([]int, length)
	for i := range input {
		input[i] = rand.Int()
	}
	slices.Sort(input)
	inputStrs := make([]string, length)
	for i := range input {
		inputStrs[i] = fmt.Sprint(input[i])
	}
	return strings.Join(inputStrs, "\n")
}

func BenchmarkReaderIndices_Search(b *testing.B) {
	var inputSize = []int{128, 256, 65536, 65536 * 2, 65536 * 4}
	for _, size := range inputSize {
		for i := 0; i < b.N; i++ {
			b.Run("input size "+fmt.Sprint(size), func(b *testing.B) {
				b.StopTimer()
				reader, err := NewReaderIndices(strings.NewReader(generateRandomSortedInts(size)))
				if err != nil {
					b.Errorf("unexpected error creating reader: %v", err)
				}
				b.StartTimer()
				_, err = reader.Search(context.Background(), reader.values[rand.Intn(len(reader.values))])
				if err != nil && !errors.Is(err, ErrIndexNotFound) {
					b.Errorf("unexpected error %v", err)
				}
			})
		}
	}
}
