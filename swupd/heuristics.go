// Copyright 2017 Intel Corporation
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

package swupd

import "strings"

func (f *File) setConfigFromPathname() {
	// TODO: make this list configurable
	configPaths := []string{
		"/etc/",
	}

	for _, path := range configPaths {
		if strings.HasPrefix(f.Name, path) {
			f.Modifier = ModifierConfig
			return
		}
	}
}

func (f *File) setStateFromPathname() {
	// TODO: make this list configurable
	statePaths := []string{
		"/usr/src/debug",
		"/dev",
		"/home",
		"/proc",
		"/root",
		"/run",
		"/sys",
		"/tmp",
		"/var",
	}

	for _, path := range statePaths {
		// if no trailing / these are state directories that are actually shipped
		// otherwise this is a non-shipped state file and the modifier should be
		// set to state
		if f.Name == path {
			return
		} else if strings.HasPrefix(f.Name, path+"/") {
			f.Modifier = ModifierState
			return
		}
	}

	// TODO: make this list configurable
	// these are paths that are not shipped directories
	stateDirs := []string{
		"/usr/src/",
		"/lost+found",
	}

	for _, path := range stateDirs {
		if strings.HasPrefix(f.Name, path) {
			f.Modifier = ModifierState
			return
		}
	}
}

func (f *File) setBootFromPathname() {
	// TODO: make this list configurable
	bootPaths := []string{
		"/boot/",
		"/usr/lib/modules/",
		"/usr/lib/kernel/",
	}

	for _, path := range bootPaths {
		if strings.HasPrefix(f.Name, path) {
			f.Modifier = ModifierBoot
			if f.Status == StatusDeleted {
				f.Status = StatusGhosted
			}
			return
		}
	}
}

func (f *File) setModifierFromPathname() {
	// order here matters, first check for config, then state, finally boot
	// more important modifiers must happen last to overwrite earlier ones
	f.setConfigFromPathname()
	f.setStateFromPathname()
	f.setBootFromPathname()
}

func (m *Manifest) applyHeuristics() {
	for _, f := range m.Files {
		f.setModifierFromPathname()
	}
}
