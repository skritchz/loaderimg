package main

import (
	"fmt"
	"github.com/rekby/gpt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s input_gpt.bin output_drive\n", os.Args[0])
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])

	if err != nil {
		fmt.Printf("Couldn't open input file with error: %v\n", err)
		os.Exit(1)
	}

	defer f.Close()

	// Skip MBR
	var dummy [0x200]byte
	_, err = f.Read(dummy[:])

	if err != nil {
		fmt.Printf("Couldn't read correctly input file: %v\n", err)
		os.Exit(1)
	}

	t, err := gpt.ReadTable(f, 512)
	if err != nil {
		fmt.Printf("Couldn't read correctly input file: %v\n", err)
		os.Exit(1)
	}

	cp := t.CreateOtherSideTable()

	//Patch all the partitions to set the _a slots bootable

	for i, p := range t.Partitions {
		name := p.Name()
		if name[len(name)-2:] == "_a" {
			t.Partitions[i].Flags[6] = 0x7f
		} else if name[len(name)-2] == '_' {
			t.Partitions[i].Flags[6] = 0x3b
		}
	}

	newFile, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Printf("Couldn't open output with error %v\n", err)
		os.Exit(2)
	}

	fc := &FakeOffset{f: newFile}

	fc.Write(dummy[:])

	if err = t.Write(fc); err != nil {
		fmt.Printf("Couldn't write new file with error %v\n", err)
		os.Exit(2)
	}
	fOffset := fc.current

	io.Copy(fc, f)

	fc.Write(dummy[:])
	fc.offset = fc.current

	if err = cp.Write(fc); err != nil {
		fmt.Printf("Couldn't open output with error %v\n", err)
		os.Exit(2)
	}
	f.Seek(fOffset, 0)
	io.Copy(fc, f)

	for _, p := range t.Partitions[:18] {
		blFile, err := os.Open(fmt.Sprintf("bootloader/%s.mbn", getName(p.Name())))
		if err != nil {
			fmt.Printf("Couldn't write sbl1:%v\n", err)
			os.Exit(3)
		}
		defer blFile.Close()

		newFile.Seek(int64(p.FirstLBA*t.SectorSize), 0)
		io.Copy(newFile, blFile)
	}
}

type FakeOffset struct {
	f       *os.File
	current int64
	offset  int64
}

func (f *FakeOffset) Seek(offset int64, whence int) (int64, error) {
	if whence == 0 {
		offset += f.offset
	}
	ret, err := f.f.Seek(offset, whence)
	if whence == 0 {
		ret -= f.offset
	}
	return ret, err

}
func (f *FakeOffset) Write(p []byte) (n int, err error) {
	w, err := f.f.Write(p)
	f.current += int64(w)
	return w, err
}

func getName(parName string) string {

	switch parName {
	case "sbl1_a", "sbl1_b":
		return "sbl1"
	case "rpm_a", "rpm_b":
		return "rpm"
	case "tz_a", "tz_b":
		return "tz"
	case "devcfg_a", "devcfg_b":
		return "devcfg"
	case "aboot_a", "aboot_b":
		return "emmc_appsboot"
	case "cmnlib_a","cmnlib_b":
		return "cmnlib_30"
	case "cmnlib64_a","cmnlib64_b":
		return "cmnlib64_30"
	case "keymaster_a", "keymaster_b":
		return "keymaster64"
	case "prov_a", "prov_b":
		return "prov"
	default:
		panic(fmt.Sprintf("unknown partition: %s", parName))
	}
}
