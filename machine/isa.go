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
	"encoding/binary"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
	"sync/atomic"
)

var isaSequence int64

// ISAYml is a YAML representation of an Instruction Set Architecture (ISA)
type ISAYml struct {
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Vendor       string   `yaml:"vendor-specific"`
	Inherits     []string `yaml:"inherits"`
	Registers    []string `yaml:"registers"`
	Instructions []string `yaml:"instructions"`
}

// ReadISAYml deserializes an ArchYml from a []byte
func ReadISAYml(raw []byte) (iy *ISAYml, err error) {
	iy = &ISAYml{}
	err = yaml.Unmarshal(raw, iy)
	return
}

// ISA is an Instruction Set Architecture (ISA) for storage in Bolt
type ISA struct {
	id           []byte
	arch         *Arch
	Name         string
	Description  string
	instructions []string
	registers    []string
	tx           *bolt.Tx
	Vendor       string
}

// ToISA converts an ISAYml to an ISA
func (i *ISAYml) ToISA(tx *bolt.Tx, a *Arch) *ISA {
	isa := &ISA{}
	id := atomic.AddInt64(&isaSequence, 1)
	buf := bytes.NewBuffer(make([]byte, 0))
	err := binary.Write(buf, binary.LittleEndian, id)
	if err != nil {
		panic(err.Error())
	}
	isa.id = buf.Bytes()
	isa.arch = a
	isa.Name = i.Name
	isa.Description = i.Description
	isa.instructions = i.Instructions
	isa.registers = i.Registers
	isa.tx = tx
	isa.Vendor = i.Vendor
	return isa
}

// Put saves an ISA to a BoltDB
func (i *ISA) Put() (err error) {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(i)
	if err != nil {
		return
	}

	b := i.tx.Bucket([]byte("isas"))
	if b == nil {
		err = errors.New("Bucket 'isas' does not exist")
		return
	}
	err = b.Put(i.id, buf.Bytes())
	if err != nil {
		return
	}

	for _, r := range i.registers {
		err = i.arch.AddReg(r, i.id)
		if err != nil {
			return
		}
	}

	for _, inst := range i.instructions {
		err = i.arch.AddInst(inst, i.id)
		if err != nil {
			return
		}
	}
	return
}
