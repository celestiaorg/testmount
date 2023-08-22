package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"

	carv2 "github.com/ipld/go-car/v2"
	"github.com/ipld/go-car/v2/blockstore"
)

const src = "DFD7A65CE5EAA57214A0F5F4C976445F34B6E333420E9423A9AEABEAC8A83C1F"
const dst = "DFD7A65CE5EAA57214A0F5F4C976445F34B6E333420E9423A9AEABEAC8A83C1F.copy"

const cidPrintCount = 5

func openCAR(path string) {
	// Open a new ReadOnly blockstore from a CARv1 file.
	// Note, `OpenReadOnly` accepts bot CARv1 and CARv2 formats and transparently generate index
	// in the background if necessary.
	// This instance sets ZeroLengthSectionAsEOF option to treat zero sized sections in file as EOF.
	robs, err := blockstore.OpenReadOnly(path,
		blockstore.UseWholeCIDs(true),
		carv2.ZeroLengthSectionAsEOF(true),
	)
	if err != nil {
		panic(err)
	}
	defer robs.Close()

	// Print root CIDs.
	roots, err := robs.Roots()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contains %v root CID(s):\n", len(roots))
	for _, r := range roots {
		size, err := robs.GetSize(context.TODO(), r)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\t%v -> %v\n", r, size)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Print the raw data size for the first 5 CIDs in the CAR file.
	keysChan, err := robs.AllKeysChan(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("List of first %v CIDs and their raw data size:\n", cidPrintCount)
	i := 1
	for k := range keysChan {
		if i > cidPrintCount {
			cancel()
			break
		}
		size, err := robs.GetSize(context.TODO(), k)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\t%v -> %v bytes\n", k, size)
		i++
	}
}

func LdWrite(w io.Writer, d ...[]byte) error {
	var sum uint64
	for _, s := range d {
		sum += uint64(len(s))
	}

	buf := make([]byte, 8)
	n := binary.PutUvarint(buf, sum)
	_, err := w.Write(buf[:n])
	if err != nil {
		return err
	}

	for _, s := range d {
		_, err = w.Write(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func Put(data []byte) error {
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	err = LdWrite(f, data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func main() {
    data, err := os.ReadFile(src)
	if err != nil {
		fmt.Println("Error reading file:", err)
		panic(err)
	}

	openCAR(src)

    var wg sync.WaitGroup
    for i:=0; i < 2; i++ {
        wg.Add(1)
        go func() {
            if err := Put(data); err != nil {
                fmt.Println(err)
            }
            wg.Done()
        }()
    }
    wg.Wait()

	openCAR(dst)
}
