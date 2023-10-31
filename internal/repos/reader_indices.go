package repos

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strconv"
)

var _ Indices = (*ReaderIndices)(nil)

// divisor for counting maximal allowed deviation from searching val
// if difference between val and found index is greater than val/[maxDeviationDiv]
// then value is not found and error is returned
const maxDeviationDiv = 10

// ReaderIndices is implementation of [Indices] repository,
// which parses indices from [io.Reader] and stores them in memory.
type ReaderIndices struct {
	values []int64
}

func NewReaderIndices(r io.Reader) (ReaderIndices, error) {
	values := make([]int64, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		curr, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return ReaderIndices{}, errors.Join(ErrInvalidIndicesFile, err)
		}
		values = append(values, curr)
	}
	if err := scanner.Err(); err != nil {
		return ReaderIndices{}, errors.Join(ErrParsingIndicesFile, err)
	}
	return ReaderIndices{values}, nil

}

// Search does binary search of given val in values of ReaderIndices.
// If val is not found, then last left and right must be the closest values to val.
// If difference between val and found index is greater than val/[maxDeviationDiv]
// then value is not found and [ErrIndexNotFound] is returned. Else index is returned.
func (ri ReaderIndices) Search(ctx context.Context, val int64) (int, error) {
	if len(ri.values) == 0 {
		return -1, ErrIndexNotFound
	}
	l, r := 0, len(ri.values)-1
	for l <= r {
		select {
		case <-ctx.Done():
			return -1, ctx.Err()
		default:
			mid := (l + r) / 2
			if ri.values[mid] == val {
				return mid, nil
			} else if ri.values[mid] > val {
				r = mid - 1
			} else {
				l = mid + 1
			}
		}

	}

	l, r = intMin(l, len(ri.values)-1), intMin(r, len(ri.values)-1)
	maxDeviation := val / maxDeviationDiv

	if intAbs(val-ri.values[r]) <= maxDeviation {
		return r, nil
	} else if intAbs(val-ri.values[l]) <= maxDeviation {
		return l, nil
	} else {
		return -1, ErrIndexNotFound
	}
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func intAbs(val int64) int64 {
	if val < 0 {
		return -val
	}
	return val
}
