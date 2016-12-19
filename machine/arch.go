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
	"encoding/gob"
	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
	"sync/atomic"
)

var archSequence int32 = -1

// ArchYml is the YAML representation of an Architecture
type ArchYml struct {
	Name  string      `yaml:"name"`
	MType elf.Machine `yaml:"type"`
}

// ReadArchYml deserializes an ArchYml from a []byte
func ReadArchYml(raw []byte) (ay *ArchYml, err error) {
	ay = &ArchYml{}
	err = yaml.Unmarshall(raw, ay)
}

// ToArch converts an ArchYml to an Arch and assigns an ID
func (ay *ArchYml) ToArch(tx *bolt.Tx) *Arch {
	a := &Arch{}
	a.id = atomic.AddInt32(&ArchSeqence, 1)
	a.Name = ay.Name
	a.mtype = ay.MType
}

// Arch represents an Architecture, with its instructions and registers
type Arch struct {
	id    int32
	Name  string
	insts *bolt.Bucket
	regs  *bolt.Bucket
	tx    *bolt.Tx
	mtype elf.Machine
}

// ReadArch deserializes an Arch from a BoltDB and retrieves its Buckets
func ReadArch(tx *bolt.Tx, aID int32) (a *Arch, err error) {
	b, err := tx.Bucket([]byte{"arch"})
	if err != nil {
		return
	}

	as, err := b.Get([]byte{strconv.Itoa(aID)})
	if err != nil {
		return
	}

	dec := gob.NewDecoder(a)
	err = dec.Decode(&as)
	if err != nil {
		return
	}

	a.id = aID
	b, err := tx.Bucket([]byte{"insts"})
	if err != nil {
		return
	}
	a.insts, err = b.Bucket([]byte{strconv.Itoa(a.id)})
	if err != nil {
		return
	}
	b, err := tx.Bucket([]byte{"regs"})
	if err != nil {
		return
	}
	a.regs, err = b.Bucket([]byte{strconv.Itoa(a.id)})
	if err != nil {
		return
	}
	a.tx = tx
}

// Put serializes an Arch into a BoltDB and creates its Buckets
func (a *Arch) Put() error {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(a)
	if err != nil {
		return
	}

	b, err := a.tx.Bucket([]byte{"insts"})
	if err != nil {
		return
	}
	a.insts, err = b.CreateBucket([]byte{strconv.Itoa(a.id)})
	if err != nil {
		return
	}

	b, err = a.tx.Bucket([]byte{"regs"})
	if err != nil {
		return
	}
	a.regs, err = b.CreateBucket([]byte{strconv.Itoa(a.id)})
	if err != nil {
		return
	}
}

// AddInst adds an instruction to an Arch
func (a *Arch) AddInst(iname string, isaID int32) error {
	err = a.insts.Put([]byte{iname}, []byte{strconv.Itoa(isaID)})
}

// InstToISA retrieves the ID for an Instruction's ISA
func (a *Arch) InstToISA(iname string) (id int32, err error) {
	v, err := a.insts.Get([]byte{iname})
	if err != nil {
		return
	}
	id, err = strconv.Atoi(v)
}

// AddReg adds a register to an Arch
func (a *Arch) AddReg(rname string, isaID int32) error {
	err = a.regs.Put([]byte{rname}, []byte{strconv.Itoa(isaID)})
}

// RegToISA retrieves the ID for a Register's ISA
func (a *Arch) RegToISA(rname string) (id int32, err error) {
	v, err := a.regs.Get([]byte{rname})
	if err != nil {
		return
	}
	id, err = strconv.Atoi(v)
}
