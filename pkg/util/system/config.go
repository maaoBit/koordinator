/*
Copyright 2022 The Koordinator Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package system

import (
	"flag"
	"os"
)

const (
	DS_MODE   = "dsMode"
	HOST_MODE = "hostMode"
)

var Conf = NewDsModeConfig()
var AgentMode = DS_MODE

type Config struct {
	CgroupRootDir         string
	CgroupKubePath        string
	SysRootDir            string
	SysFSRootDir          string
	ProcRootDir           string
	VarRunRootDir         string
	NodeNameOverride      string
	RuntimeHooksConfigDir string

	ContainerdEndPoint string
	DockerEndPoint     string

	UseCgroupsV2 bool // TODO: guess the cgroup version according to cgroup path
}

func NewHostModeConfig() *Config {
	return &Config{
		CgroupKubePath:        "kubepods/",
		CgroupRootDir:         "/sys/fs/cgroup/",
		ProcRootDir:           "/proc/",
		SysRootDir:            "/sys/",
		SysFSRootDir:          "/sys/fs/",
		VarRunRootDir:         "/var/run/",
		RuntimeHooksConfigDir: "/etc/runtime/hookserver.d",
		UseCgroupsV2:          false,
	}
}

func NewDsModeConfig() *Config {
	return &Config{
		CgroupKubePath: "kubepods/",
		CgroupRootDir:  "/host-cgroup/",
		// some dirs are not covered by ns, or unused with `hostPID` is on
		ProcRootDir:           "/proc/",
		SysRootDir:            "/host-sys/",
		SysFSRootDir:          "/host-sys-fs/",
		VarRunRootDir:         "/host-var-run/",
		RuntimeHooksConfigDir: "/host-etc-hookserver/",
		UseCgroupsV2:          false,
	}
}

func init() {
	agentMode := os.Getenv("agent_mode")
	if agentMode == HOST_MODE {
		Conf = NewHostModeConfig()
		AgentMode = agentMode
	}
	initFilePath()
}

func SetConf(config Config) {
	Conf = &config
	HostSystemInfo = collectVersionInfo()
	initFilePath()
}

func (c *Config) InitFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.CgroupRootDir, "cgroup-root-dir", c.CgroupRootDir, "Cgroup root dir")
	fs.StringVar(&c.SysRootDir, "sys-root-dir", c.SysRootDir, "host /sys dir in container")
	fs.StringVar(&c.SysFSRootDir, "sys-fs-root-dir", c.SysFSRootDir, "host /sys/fs dir in container, used by resctrl fs")
	fs.StringVar(&c.ProcRootDir, "proc-root-dir", c.ProcRootDir, "host /proc dir in container")
	fs.StringVar(&c.VarRunRootDir, "var-run-root-dir", c.VarRunRootDir, "host /var/run dir in container")

	fs.StringVar(&c.CgroupKubePath, "cgroup-kube-dir", c.CgroupKubePath, "Cgroup kube dir")
	fs.StringVar(&c.NodeNameOverride, "node-name-override", c.NodeNameOverride, "If non-empty, will use this string as identification instead of the actual machine name. ")
	fs.StringVar(&c.ContainerdEndPoint, "containerd-endpoint", c.ContainerdEndPoint, "containerd endPoint")
	fs.StringVar(&c.DockerEndPoint, "docker-endpoint", c.DockerEndPoint, "docker endPoint")

	fs.BoolVar(&c.UseCgroupsV2, "use-cgroups-v2", c.UseCgroupsV2, "whether use cgroups-v2 or not")

	HostSystemInfo = collectVersionInfo()
	initFilePath()
}
