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

var instMatch = regexp.MustCompile("^\\s+\\w+:\\s+(\\w+)")
var regMatch = regexp.MustCompile("%(\\w+)")

// ParseLine gets the instruction and registers from a single line of objdump output
func ParseLine(line []byte) (inst string, regs []string, ok bool) {
	ok = false
	regs = make([]string, 0)
	l := string(line[:])
	// find instruction
	m := instMatch.FindStringSubmatch(l)
	if len(m) != 2 {
		return
	}
	ok = true
	inst = m[1]
	ms := regMatch.FindAllStringSubmatch(l, -1)
	for _, m := range ms {
		regs = append(regs, m[1])
	}
	return
}

// ReadObjdump gets the counts of all registers and instructions used in a binary file
func ReadObjdump(stdout io.Reader) (insts map[string]int64, regs map[string]int64, err error) {
	insts = make(map[string]int64)
	regs = make(map[string]int64)
	r := bufio.NewReaderSize(stdout, 100)
	for {
		line, e := r.ReadBytes('\n')
		if e != nil {
			if e == io.EOF {
				break
			}
			err = e
			return
		}
		i, rs, ok := ParseLine(line)
		if ok {
			insts[i]++
			for _, r := range rs {
				regs[r]++
			}
		}
	}
	return
}

// RunObjdump executes objdump, getting the counts of all registers and instructions used in the specified file
func RunObjdump(fpath string) (insts map[string]int64, regs map[string]int64, err error) {
	cmd := exec.Command("objdump", "--no-show-raw-insn", "-d", fpath)
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
