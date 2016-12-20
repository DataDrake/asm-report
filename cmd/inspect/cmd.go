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
	"github.com/DataDrake/asm-report/machine"
	"github.com/boltdb/bolt"
	"os"
)

func usage() {
	print("USAGE: asm-report inspect [OPTIONS] FILE...\n")
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
		println(arch.Name)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
