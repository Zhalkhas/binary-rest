package repos

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strconv"
)

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

func (f ReaderIndices) Search(ctx context.Context, val int) (int, error) {
	//TODO implement me
	panic("implement me")
}
