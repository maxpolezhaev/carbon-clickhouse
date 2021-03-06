package uploader

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/lomik/carbon-clickhouse/helper/RowBinary"
)

type Series struct {
	*cached
	isReverse bool
}

var _ Uploader = &Series{}
var _ UploaderWithReset = &Series{}

func NewSeries(base *Base) *Series {
	u := &Series{}
	u.cached = newCached(base)
	u.cached.parser = u.parseFile
	return u
}

func NewSeriesReverse(base *Base) *Series {
	u := NewSeries(base)
	u.isReverse = true
	return u
}

func (u *Series) parseFile(filename string, out io.Writer) (map[string]bool, error) {
	var reader *RowBinary.Reader
	var err error

	if u.isReverse {
		reader, err = RowBinary.NewReverseReader(filename)
	} else {
		reader, err = RowBinary.NewReader(filename)
	}

	if err != nil {
		return nil, err
	}
	defer reader.Close()

	version := uint32(time.Now().Unix())

	newSeries := make(map[string]bool)

	wb := RowBinary.GetWriteBuffer()

	var level int

LineLoop:
	for {
		name, err := reader.ReadRecord()
		if err != nil { // io.EOF or corrupted file
			break
		}

		// skip tagged
		if bytes.IndexByte(name, '?') >= 0 {
			continue
		}

		key := fmt.Sprintf("%d:%s", reader.Days(), unsafeString(name))

		if u.existsCache.Exists(key) {
			continue LineLoop
		}

		if newSeries[key] {
			continue LineLoop
		}

		level = pathLevel(name)

		wb.Reset()

		newSeries[key] = true
		wb.WriteUint16(reader.Days())
		wb.WriteUint32(uint32(level))
		wb.WriteBytes(name)
		wb.WriteUint32(version)

		_, err = out.Write(wb.Bytes())
		if err != nil {
			return nil, err
		}
	}

	wb.Release()

	return newSeries, nil
}
