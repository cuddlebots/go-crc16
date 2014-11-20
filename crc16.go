// Copyright 2009 The Go Authors. All rights reserved.
// Copyright 2014 Michael Phan-Ba <michael@mikepb.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crc16 implements the 16-bit cyclic redundancy check, or CRC-16,
// checksum. See http://en.wikipedia.org/wiki/Cyclic_redundancy_check for
// information.
package crc16

// The size of a CRC-16 checksum in bytes.
const Size = 2

// https://en.wikipedia.org/wiki/Cyclic_redundancy_check#Standards_and_common_use
const (
	// Bisync, Modbus, USB, ANSI X3.28, SIA DC-07, many others
	ANSI = 0x8005
	// X.25, V.41, HDLC FCS, XMODEM, Bluetooth, PACTOR, SD, many others
	CCITT = 0x1021
)

type Table [256]uint16

// ANSITable is the table for the ANSI polynomial.
var ANSITable = makeTable(ANSI)

// CCITTTable is the table for the CCITT polynomial.
var CCITTTable = makeTable(CCITT)

// MakeTable returns the Table constructed from the specified polynomial.
func MakeTable(poly uint16) *Table {
	switch poly {
	case ANSI:
		return ANSITable
	case CCITT:
		return CCITTTable
	}
	return makeTable(poly)
}

// makeTable returns the Table constructed from the specified polynomial.
func makeTable(poly uint16) *Table {
	t := new(Table)
	for i := 0; i < 256; i++ {
		crc := uint16(i)
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ poly
			} else {
				crc >>= 1
			}
		}
		t[i] = crc
	}
	return t
}

// digest represents the partial evaluation of a checksum.
type digest struct {
	crc uint16
	tab *Table
}

// New creates a new Hash16 computing the CRC-32 checksum
// using the polynomial represented by the Table.
func New(tab *Table) Hash16 { return &digest{0, tab} }

// NewANSI creates a new Hash16 computing the CRC-32 checksum
// using the ANSI polynomial.
func NewANSI() Hash16 { return New(ANSITable) }

// NewCCITT creates a new Hash16 computing the CRC-32 checksum
// using the CCITT polynomial.
func NewCCITT() Hash16 { return New(CCITTTable) }

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return 1 }

func (d *digest) Reset() { d.crc = 0 }

// Update returns the result of adding the bytes in p to the crc.
func Update(crc uint16, tab *Table, p []byte) uint16 {
	crc = ^crc
	for _, v := range p {
		crc = tab[byte(crc)^v] ^ (crc >> 8)
	}
	return ^crc
}

func (d *digest) Write(p []byte) (n int, err error) {
	d.crc = Update(d.crc, d.tab, p)
	return len(p), nil
}

func (d *digest) Sum16() uint16 { return d.crc }

func (d *digest) Sum(in []byte) []byte {
	s := d.Sum16()
	return append(in, byte(s>>8), byte(s))
}

// Checksum returns the CRC-16 checksum of data
// using the polynomial represented by the Table.
func Checksum(data []byte, tab *Table) uint16 { return Update(0, tab, data) }

// ChecksumANSI returns the CRC-16 checksum of data
// using the ANSI polynomial.
func ChecksumANSI(data []byte) uint16 { return Update(0, ANSITable, data) }

// ChecksumCCITT returns the CRC-16 checksum of data
// using the CCITT polynomial.
func ChecksumCCITT(data []byte) uint16 { return Update(0, CCITTTable, data) }
