package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/filecoin-project/dagstore/mount"
)

type fileMount struct {
	buf *bytes.Buffer
	mount.FileMount
}

func (m *fileMount) WriteTo(w io.Writer) (int64, error) {
	return io.Copy(w, bytes.NewReader(m.buf.Bytes()))
}

func Write(data []byte, mount *fileMount) error {
    _, err := mount.buf.Write(data)
    return err
}

func Put() error {
    key := "block001.dat"
	f, err := os.OpenFile(key, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// save encoded eds into buffer
	m := &fileMount{
		// TODO: buffer could be pre-allocated with capacity calculated based on eds size.
		buf:       bytes.NewBuffer(nil),
		FileMount: mount.FileMount{Path: key},
	}
	randomBytes := make([]byte, 1024 * 4096)
	_, err = rand.Read(randomBytes)
	if err != nil {
		fmt.Println("Error generating random bytes:", err)
		return err
	}
	err = Write(randomBytes, m)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	// write whole buffered mount data in one go to optimize i/o
	if _, err = m.WriteTo(f); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func main() {
    var wg sync.WaitGroup
    for i:=0; i < 2; i++ {
        wg.Add(1)
        go func() {
            if err := Put(); err != nil {
                fmt.Println(err)
            }
            wg.Done()
        }()
    }
    wg.Wait()
}
