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

package inspect

import (
	"debug/elf"
	"fmt"
	"github.com/DataDrake/asm-report/machine"
	"github.com/boltdb/bolt"
	"os"
	"sort"
)

func usage() {
	print("USAGE: asm-report inspect [OPTIONS] FILE...\n")
}

func getISAs(tx *bolt.Tx, arch *machine.Arch, insts, regs map[string]int64) (isas map[string]int64, err error) {
	isas = make(map[string]int64)
	for k, v := range insts {
		isaID, e := arch.InstToISA(k)
		if e != nil {
			println(k)
			continue
		}
		isa, e := machine.ReadISA(tx, isaID)
		if e != nil {
			println(k)
			err = e
			return
		}
		isas[isa.Name] += v
	}

	for k, v := range regs {
		isaID, e := arch.RegToISA(k)
		if e != nil {
			println(k)
			continue
		}
		isa, e := machine.ReadISA(tx, isaID)
		if e != nil {
			err = e
			return
		}
		isas[isa.Name] += v
	}
	return
}

func printMap(m map[string]int64) {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("    - %-12s : %d\n", key, m[key])
	}
}

// Cmd handles the "inspect" subcommand
func Cmd(args []string) {
	db, err := bolt.Open("/tmp/test.db", 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}
	f, err := elf.Open(args[0])
	if err != nil {
		panic(err.Error())
	}
	mtype := f.FileHeader.Machine
	f.Close()
	err = db.View(func(tx *bolt.Tx) error {
		arch, err := machine.ReadArchElf(tx, mtype)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%s : %s\n", "Architecture", arch.Name)
		insts, regs, err := machine.RunObjdump(args[0])
		if err != nil {
			return err
		}
		//fmt.Printf("%s : %d\n", "Instructions", len(insts))
		//printMap(insts)
		//fmt.Printf("%s : %d\n", "Registers", len(regs))
		//printMap(regs)
		isas, err := getISAs(tx, arch, insts, regs)
		if err != nil {
			return err
		}
		fmt.Printf("%s : %d\n", "ISAS", len(isas))
		printMap(isas)
		return err
	})
	if err != nil {
		panic(err.Error())
	}
}
