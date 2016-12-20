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

package info

import (
    "fmt"
	"github.com/DataDrake/asm-report/machine"
	"github.com/boltdb/bolt"
)

func printArch(a *machine.Arch, full bool) {
    if full {
        fmt.Printf("%-12s : %s\n", "Architecture", a.Name)
        fmt.Printf("%-12s : %d\n", "Instructions", a.NInst())
        fmt.Printf("%-12s : %d\n", "Registers", a.NReg())
    } else {
        fmt.Printf("    - %s\n", a.Name)
    }
}

func usage() {
	print("USAGE: asm-report info [OPTIONS]\n")
}

func Cmd(args []string) {
	db, err := bolt.Open("/tmp/test.db", 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
    if len(args) == 0 {
        fmt.Println("Available architectures:")
    }
    db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("arch"))
        c := b.Cursor()
        for k,_ := c.First(); k != nil; k,_ = c.Next() {
            a, err := machine.ReadArch(tx, k)
            if err != nil {
                return err
            }
        	if len(args) > 0 {
                if a.Name == args[0] {
                    printArch(a, true)
                }
        	} else {
                printArch(a, false)
            }
        }
        return nil
    })
}
