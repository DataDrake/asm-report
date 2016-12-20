//
// Copyright © 2016 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implie$
// See the License for the specific language governing permissions and
// limitations under the License.
//

package machine

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
	"sync/atomic"
)

var archSequence int64

// ArchYml is the YAML representation of an Architecture
type ArchYml struct {
	Name   string        `yaml:"name"`
	MTypes []elf.Machine `yaml:"types"`
}

// ReadArchYml deserializes an ArchYml from a []byte
func ReadArchYml(raw []byte) (ay *ArchYml, err error) {
	ay = &ArchYml{}
	err = yaml.Unmarshal(raw, ay)
	return
}

// ToArch converts an ArchYml to an Arch and assigns an ID
func (ay *ArchYml) ToArch(tx *bolt.Tx) *Arch {
	a := &Arch{}
	id := atomic.AddInt64(&archSequence, 1)
	buf := bytes.NewBuffer(make([]byte, 0))
	err := binary.Write(buf, binary.LittleEndian, id)
	if err != nil {
		panic(err.Error())
	}
	a.id = buf.Bytes()
	a.Name = ay.Name
	a.mtypes = ay.MTypes
	a.tx = tx
	return a
}

// Arch represents an Architecture, with its instructions and registers
type Arch struct {
	id     []byte
	Name   string
	insts  *bolt.Bucket
	regs   *bolt.Bucket
	tx     *bolt.Tx
	mtypes []elf.Machine
}

// ReadArch deserializes an Arch from a BoltDB and retrieves its Buckets
func ReadArch(tx *bolt.Tx, id []byte) (a *Arch, err error) {
	b := tx.Bucket([]byte("arch"))
	if b == nil {
		err = errors.New("Bucket 'arch' does not exist")
		return
	}

	as := b.Get(id)
	if as == nil {
		err = errors.New("Architecture not found")
		return
	}

	bbuf := bytes.NewBuffer(as)
	dec := gob.NewDecoder(bbuf)
	err = dec.Decode(&a)
	if err != nil {
		return
	}

	a.id = id
	b = tx.Bucket([]byte("insts"))
	if b == nil {
		err = errors.New("Bucket 'insts' does not exist")
		return
	}
	a.insts = b.Bucket(id)
	if a.insts == nil {
		err = errors.New("'insts' arch-specific subbucket does not exist")
		return
	}
	b = tx.Bucket([]byte("regs"))
	if b == nil {
		err = errors.New("Bucket 'regs' does not exist")
		return
	}
	a.regs = b.Bucket(id)
	if a.regs == nil {
		err = errors.New("'regs' arch-specific subbucket does not exist")
		return
	}
	a.tx = tx
	return
}

// Put serializes an Arch into a BoltDB and creates its Buckets
func (a *Arch) Put() (err error) {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(a)
	if err != nil {
		return
	}

	b := a.tx.Bucket([]byte("arch"))
	if b == nil {
		err = errors.New("Bucket 'arch' does not exist")
		return
	}
	err = b.Put(a.id, buf.Bytes())
	if err != nil {
		return
	}
	b = a.tx.Bucket([]byte("elf"))
	if b == nil {
		err = errors.New("Bucket 'elfs' does not exist")
		return
	}
	for _, e := range a.mtypes {
		buf := bytes.NewBuffer(make([]byte, 0))
		err = binary.Write(buf, binary.LittleEndian, e)
		if err != nil {
			return
		}
		err = b.Put(buf.Bytes(), a.id)
		if err != nil {
			return
		}
	}

	b = a.tx.Bucket([]byte("insts"))
	if b == nil {
		err = errors.New("Bucket 'insts' does not exist")
		return
	}
	a.insts, err = b.CreateBucket(a.id)
	if err != nil {
		return
	}

	b = a.tx.Bucket([]byte("regs"))
	if b == nil {
		err = errors.New("Bucket 'regs' does not exist")
		return
	}
	a.regs, err = b.CreateBucket(a.id)
	if err != nil {
		return
	}
	return
}

// AddInst adds an instruction to an Arch
func (a *Arch) AddInst(iname string, isaID []byte) (err error) {
	err = a.insts.Put([]byte(iname), isaID)
	return
}

// InstToISA retrieves the ID for an Instruction's ISA
func (a *Arch) InstToISA(iname string) (id []byte, err error) {
	id = a.insts.Get([]byte(iname))
	if id == nil {
		err = errors.New("Instruction not found")
		return
	}
	return
}

// NInst gets the total number of instructions for an Arch
func (a *Arch) NInst() int {
    return a.insts.Stats().KeyN
}

// AddReg adds a register to an Arch
func (a *Arch) AddReg(rname string, isaID []byte) (err error) {
	err = a.regs.Put([]byte(rname), isaID)
	return
}

// RegToISA retrieves the ID for a Register's ISA
func (a *Arch) RegToISA(rname string) (id []byte, err error) {
	v := a.regs.Get([]byte(rname))
	if v == nil {
		err = errors.New("Register not found")
		return
	}
	return
}

// NReg gets the total number of registers for an Arch
func (a *Arch) NReg() int {
    return a.regs.Stats().KeyN
}