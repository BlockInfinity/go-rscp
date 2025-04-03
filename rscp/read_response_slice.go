package rscp

import (
	"fmt"

	"github.com/spali/go-slicereader"
)

// readResponseSlice reads a single response message (recursive for containers) from a slice.
func readResponseSlice(values []interface{}) (*Message, error) {
	sr := slicereader.NewSliceReader(values)
	return readResponseSliceReader(sr)
}

// readResponseSliceReader reads a single response message (recursive for containers) from the SliceReader.
func readResponseSliceReader(sr *slicereader.SliceReader) (*Message, error) {
	var (
		t     interface{}
		tag   Tag
		isTag bool
		err   error
	)
	if t, err = sr.Read(); err != nil {
		return nil, err
	}
	if tag, isTag = t.(Tag); !isTag {
		return nil, fmt.Errorf("element at index %d: %w", int(sr.Size()-int64(sr.Len())-1), ErrValidTag)
	}
	msg := NewMessage(tag, nil)
	if msg.DataType == None {
		return msg, nil
	}
	if msg.DataType == Container {
		msg.Value = make([]Message, 0)
		for sr.Len() > 0 {
			var (
				subMsg *Message
				err    error
			)
			if subMsg, err = readResponseSliceReader(sr); err != nil {
				return nil, err
			}
			msg.Value = append(msg.Value.([]Message), *subMsg)
		}
	} else {
		if msg.Value, err = sr.Read(); err != nil {
			return nil, fmt.Errorf("expect value after tag %s with data type %s: %w", msg.Tag, msg.DataType, ErrMissingValue)
		}
		switch msg.Value.(type) {
		case Tag, DataType:
			return nil, fmt.Errorf("expect element at index %d to be a value for tag %s, got type: %T: %w",
				int(sr.Size()-int64(sr.Len())-1), msg.Tag, msg.Value, ErrDataTypeValueMismatch)
		}
	}
	return msg, nil
}
