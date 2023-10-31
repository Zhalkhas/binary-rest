package repos

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
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
	slog.Debug("start parsing indices from reader")
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		curr, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			slog.Error("error while parsing string line to int", "err", err, "line", scanner.Text())
			return ReaderIndices{}, errors.Join(ErrInvalidIndicesFile, err)
		}
		values = append(values, curr)
	}
	if err := scanner.Err(); err != nil {
		slog.Error("error while parsing indices from reader", "err", err)
		return ReaderIndices{}, errors.Join(ErrParsingIndicesFile, err)
	}
	slog.Debug("end parsing indices from reader, parsed values count:" + fmt.Sprint(len(values)))
	return ReaderIndices{values}, nil

}

// Search does binary search of given val in values of ReaderIndices.
// If val is not found, then last left and right must be the closest values to val.
// If difference between val and found index is greater than val/[maxDeviationDiv]
// then value is not found and [ErrIndexNotFound] is returned. Else index is returned.
func (ri ReaderIndices) Search(ctx context.Context, val int64) (int, error) {
	slog.Debug("start searching index for value", "value", val)
	if len(ri.values) == 0 {
		slog.Error("values slice is empty")
		return -1, ErrIndexNotFound
	}
	l, r := 0, len(ri.values)-1
	for l <= r {
		select {
		case <-ctx.Done():
			slog.Error("context is done while searching index", "err", ctx.Err())
			return -1, ctx.Err()
		default:
			mid := (l + r) / 2
			if ri.values[mid] == val {
				slog.Debug("end searching index for value", "value", val, "index", mid)
				return mid, nil
			} else if ri.values[mid] > val {
				r = mid - 1
			} else {
				l = mid + 1
			}
		}
	}

	abs := func(val int64) int64 {
		if val < 0 {
			return -val
		}
		return val
	}

	slog.Debug("couldn't find exact value, trying to identify closest")

	l, r = min(l, len(ri.values)-1), min(r, len(ri.values)-1)
	lDeviation, rDeviation := abs(val-ri.values[l]), abs(val-ri.values[r])
	maxDeviation := val / maxDeviationDiv

	if rDeviation <= maxDeviation {
		slog.Debug("end searching index for value, right chosen", "value", val, "index", r)
		return r, nil
	} else if lDeviation <= maxDeviation {
		slog.Debug("end searching index for value, left chosen", "value", val, "index", r)
		return l, nil
	} else {
		slog.Error("couldn't find closest value", "value", val, "left", l, "right", r)
		return -1, ErrIndexNotFound
	}
}
