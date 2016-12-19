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
    "io/ioutil"
    "os"
)

func usage() {
	print("USAGE: asm-report update [OPTIONS]\n")
}

func readDefs(d *os.File) (ays []*machine.ArchYml, isas []*machine.ISAYml, err error) {
    fs, err := d.Readdir(-1)
    if err != nil {
        return
    }

    ays  = make([]*machine.ArchYml,0)
    archfs := make([]os.FileInfo,0)
    isas = make([]*machine.ISAYml,0)
    isadirs := make([]os.FileInfo,0)

    for _,f := range fs {
        if f.IsDir() {
            isadirs = append(isadirs, f)
        } else {
            archfs = append(archfs, f)
        }
    }

    for _, archf := range archfs {
        af, e := os.Open("./defs/" + archf.Name())
        defer af.Close()
        if e != nil {
            err = e
            return
        }
        ayraw, e := ioutil.ReadAll(af)
        if e != nil {
            err = e
            return
        }
        ay, e := machine.ReadArchYml(ayraw)
        if e != nil {
            err = e
            return
        }
        ays = append(ays, ay)
    }
    println(len(ays))
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
    _, _, err = readDefs(d)
}
