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

package update

import (
	"github.com/DataDrake/asm-report/machine"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"os"
)

func usage() {
	print("USAGE: asm-report update [OPTIONS]\n")
}

func updateISA(tx *bolt.Tx, arch *machine.Arch, iname string) {
	ifile, err := os.Open("./defs/" + arch.Name + "/" + iname)
	if err != nil {
		panic(err.Error())
	}
	defer ifile.Close()
	iraw, err := ioutil.ReadAll(ifile)
	if err != nil {
		panic(err.Error())
	}
	iy, err := machine.ReadISAYml(iraw)
	if err != nil {
		println(iname)
		panic(err.Error())
	}
	i := iy.ToISA(tx, arch)
	err = i.Put()
	if err != nil {
		panic(err.Error())
	}
}

func updateISAs(tx *bolt.Tx, arch *machine.Arch) (err error) {
	d, err := os.Open("./defs/" + arch.Name)
	if err != nil {
		return
	}
	defer d.Close()
	isafs, err := d.Readdir(-1)
	if err != nil {
		return
	}
	for _, isaf := range isafs {
		if !isaf.IsDir() {
			updateISA(tx, arch, isaf.Name())
		}
	}
	return
}

func updateArch(db *bolt.DB, aname string) (err error) {
	af, err := os.Open("./defs/" + aname + ".yml")
	if err != nil {
		return
	}
	defer af.Close()
	ayraw, err := ioutil.ReadAll(af)
	if err != nil {
		return
	}
	ay, err := machine.ReadArchYml(ayraw)
	if err != nil {
		return
	}
	err = db.Update(func(tx *bolt.Tx) error {
		a := ay.ToArch(tx)
		err := a.Put()
		if err != nil {
			return err
		}
		err = updateISAs(tx, a)
		return err
	})
	return
}

func updateDefs(d *os.File) (err error) {
	fs, err := d.Readdir(-1)
	if err != nil {
		return
	}
	db, err := bolt.Open("/tmp/test2.db", 0600, nil)
	if err != nil {
		return
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		var e error
		_, e = tx.CreateBucket([]byte("arch"))
		if e != nil {
			return e
		}
		_, e = tx.CreateBucket([]byte("elf"))
		if e != nil {
			return e
		}
		_, e = tx.CreateBucket([]byte("insts"))
		if e != nil {
			return e
		}
		_, e = tx.CreateBucket([]byte("isas"))
		if e != nil {
			return e
		}
		_, e = tx.CreateBucket([]byte("regs"))
		return e
	})
	if err != nil {
		return
	}

	for _, f := range fs {
		if f.IsDir() {
			e := updateArch(db, f.Name())
			if e != nil {
				err = e
				return
			}
		}
	}
	err = os.Rename("/tmp/test2.db", "/tmp/test.db")
	return
}

// Cmd handles the "update" subcommand
func Cmd(args []string) {
	if len(args) != 0 {
		usage()
		os.Exit(1)
	}

	d, err := os.Open("./defs")
	if err != nil {
		panic(err.Error())
	}
	defer d.Close()
	err = updateDefs(d)
	if err != nil {
		panic(err.Error())
	}
}
