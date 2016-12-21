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
	"regexp"
)

var instMatch = regexp.MustCompile(":\\s+(\\w+)")
var regMatch = regexp.MustCompile("%\\w+")

// ReadObjdump gets the counts of all registers and instructions used in a binary file
func ReadObjdump(stdout io.Reader) (insts map[string]int64, regs map[string]int64, err error) {
	insts = make(map[string]int64)
	regs = make(map[string]int64)
	r := bufio.NewReaderSize(stdout, 100)
    var line []byte
	for {
		//line, e := r.ReadBytes('\n')
		line, _, err = r.ReadLine()
		if err != nil {
			if err == io.EOF {
                err = nil
			    return
			}
			return
		}
       	m := instMatch.FindSubmatch(line)
        if len(m) != 2 {
            continue
        }
        insts[string(m[1])]++
        ms := regMatch.FindAll(line, -1)
        for _, m := range ms {
            regs[string(m[1:])]++
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
