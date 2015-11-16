/*
   Copyright 2014 Franc[e]sco (lolisamurai@tfwno.gf)
   This file is part of sharenix.
   sharenix is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   sharenix is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with sharenix. If not, see <http://www.gnu.org/licenses/>.
*/

package sharenixlib

import (
	"fmt"
	"github.com/kardianos/osext"
	"os/exec"
	"path"
	"strings"
)

// GetPluginsDir returns the absolute path to the plugins directory.
func GetPluginsDir() (pluginsDir string, err error) {
	exeFolder, err := osext.ExecutableFolder()
	if err != nil {
		return
	}
	pluginsDir = path.Join(exeFolder, "/plugins/")
	return
}

// RunPlugin starts pluginName in the plugin directory passing command-line
// params in the following format:
// 	pluginName -param1Name=param1Value ... -paramXName=paramXValue param_tail
// For example, calling
// 	RunPlugin("foo", map[string]string{
// 		"hello": "world",
// 		"someflag": "true",
//		"_tail": "bar",
// 	})
// will execute
// 	foo -hello=world -someflag=true bar
// Returns the last line outputted to stdout by the plugin and an error if any.
// Note: strings.Fields() is called on tail, so any extra spaces will be
// stripped.
func RunPlugin(pluginName string,
	extraParams map[string]string) (output string, err error) {

	formattedArgs := strings.Fields(extraParams["_tail"])
	delete(extraParams, "_tail")
	for paramName, paramValue := range extraParams {
		formattedArgs = append(
			[]string{fmt.Sprintf("-%s=%s", paramName, paramValue)},
			formattedArgs...)
	}

	pluginsDir, err := GetPluginsDir()
	if err != nil {
		return
	}

	outdata, err := exec.Command(path.Join(pluginsDir, pluginName),
		formattedArgs...).CombinedOutput()
	DebugPrintln("exec.CombinedOutput returned:\n",
		string(outdata), "with error", err)
	if len(outdata) == 0 && err == nil {
		err = fmt.Errorf("Plugin did not return any output.")
		return
	}
	split := strings.Split(strings.TrimSuffix(string(outdata), "\n"), "\n")
	output = split[len(split)-1]
	return
}