//
// Copyright 2017 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"github.com/DataDrake/asm-report/cmd/info"
	"github.com/DataDrake/asm-report/cmd/inspect"
	"github.com/DataDrake/asm-report/cmd/update"
	"github.com/pkg/profile"
	"os"
)

func usage() {
	print("USAGE: asm-report <CMD> [OPTIONS] [ARGS]\n")
}

func main() {
	defer profile.Start().Stop()
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "info":
		info.Cmd(os.Args[2:])
	case "inspect":
		inspect.Cmd(os.Args[2:])
	case "update", "up":
		update.Cmd(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}

	//os.Exit(0)
}
