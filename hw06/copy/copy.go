package copy

import (
	"errors"
	"io"
	"math"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var maxBufferSize = 1024

// Copy data from - to file
func Copy(from string, to string, limit int, offset int) error {
	if limit < 0 {
		return errors.New("limit mast be positive")
	}

	if offset < 0 {
		return errors.New("offset mast be positive")
	}

	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	info, err := os.Stat(from)
	if err != nil {
		return err
	}

	if int64(offset) > info.Size() {
		return errors.New("offset over file size")
	}

	if int64(limit) > info.Size()-int64(offset) {
		return errors.New("limit over fileSize - offset")
	}

	if limit == 0 {
		limit = int(info.Size() - int64(offset))
	}

	_, err = fromFile.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return err
	}

	toFile, err := os.Create(to)
	if err != nil {
		return err
	}

	defer toFile.Close()

	// create and start new bar
	bar := pb.StartNew(limit)
	defer bar.Finish()

	for i := offset; i < offset+limit; {
		bufferSize := int(math.Min(float64(offset+limit-i), float64(maxBufferSize)))
		buf := make([]byte, bufferSize)

		// read from file
		currentOffset := 0
		for currentOffset < bufferSize {
			read, err := fromFile.Read(buf[currentOffset:])
			currentOffset += read
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}

		// write to file
		_, err := toFile.Write(buf)
		if err != nil {
			return err
		}

		i += bufferSize

		bar.Add(bufferSize)
	}

	return nil
}
