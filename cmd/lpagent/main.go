// Copyright 2015 Authors, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"syscall"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"
	"github.com/hooto/httpsrv/deps/go.net/websocket"

	"code.hooto.com/lessos/loscore/config"
	"code.hooto.com/lessos/loscore/losapi"
	"code.hooto.com/lessos/loscore/lpagent/executor"
	"code.hooto.com/lessos/loscore/lpagent/v1"
)

const (
	addr_sock = "/home/action/.los/lpagent.sock"
)

var (
	pod_id = ""
)

func main() {

	runtime.GOMAXPROCS(1)

	//
	pod_id = strings.TrimSpace(os.Getenv("POD_ID"))
	if !losapi.PodIdReg.MatchString(pod_id) {
		os.Exit(1)
	}

	//
	if _, err := user.Lookup(config.User.Username); err != nil {

		nologin, err := exec.LookPath("nologin")
		if err != nil {
			nologin = "/sbin/nologin"
		}

		if _, err = exec.Command(
			"/usr/sbin/useradd",
			"-d", "/home/action",
			"-s", nologin,
			"-u", config.User.Uid, config.User.Username,
		).Output(); err != nil {
			os.Exit(1)
		}
	}

	//
	syscall.Setgid(2048)
	syscall.Setuid(2048)
	syscall.Chdir("/home/action")

	//
	os.MkdirAll("/home/action/var/log", 0755)
	hlog4g.LogDirSet("/home/action/var/log")
	hlog.Printf("info", "started")

	//
	go executor.Runner("/home/action")

	//
	httpsrv.GlobalService.Config.HttpAddr = "unix:" + addr_sock

	httpsrv.GlobalService.HandlerRegister(
		"/los/v1/pb/termws",
		websocket.Handler(v1.TerminalWsOpenAction))

	httpsrv.GlobalService.ModuleRegister("/los/v1/", v1.NewModule())

	httpsrv.GlobalService.Start()

	select {}
}
