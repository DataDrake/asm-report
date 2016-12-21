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
	"bufio"
	"io"
	"os/exec"
    "strings"
)

// ReadObjdump gets the counts of all registers and instructions used in a binary file
func ReadObjdump(stdout io.Reader) (insts map[string]int64, regs map[string]int64, err error) {
	insts = make(map[string]int64)
	regs = make(map[string]int64)
	r := bufio.NewReaderSize(stdout, 100)
    var line []byte
	for {
		line, _, err = r.ReadLine()
		if err != nil {
			if err == io.EOF {
                err = nil
			    return
			}
			return
		}
        sl := string(line)
        if strings.Index(sl," ") != 0 {
            continue
        }
        i := strings.Index(sl,":") + 1
        inst := strings.TrimSpace(sl[i:])
        i = strings.IndexAny(inst," \t")
        if i < 0 {
           insts[inst]++
        } else {
            insts[inst[:i]]++
        }
        for i,reg := range strings.Split(sl,"%") {
            if i < 1 {
                continue
            }
            j := strings.IndexAny(reg,",:)")
            if j > 0 {
                regs[reg[:j]]++
            } else {
                regs[reg]++
            }
        }
	}
	return
}

// RunObjdump executes objdump, getting the counts of all registers and instructions used in the specified file
func RunObjdump(fpath string) (insts map[string]int64, regs map[string]int64, err error) {
	//cmd := exec.Command("objdump", "--no-show-raw-insn", "-d", fpath)
	cmd := exec.Command("llvm-objdump", "-no-show-raw-insn", "-disassemble", fpath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	err = cmd.Start()
	if err != nil {
		return
	}
	insts, regs, err = ReadObjdump(stdout)
	if err != nil {
		return
	}
	err = cmd.Wait()
	return
}
