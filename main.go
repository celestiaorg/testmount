package main

import (
	"bytes"
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
	data := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum bibendum turpis imperdiet dolor faucibus aliquet. Phasellus viverra tellus a placerat dapibus. Integer vitae dui ornare, egestas velit eget, fermentum tortor. Quisqueeget nibh id orci ultrices pulvinar. In hac habitasse platea dictumst.")
	data = append(data, []byte{0x0a}...)
	err = Write(data, m)
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
    for i:=0; i < 100; i++ {
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
