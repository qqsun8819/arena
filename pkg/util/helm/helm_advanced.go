// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helm

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
* Generate value file
 */
func GenerateValueFile(values interface{}) (valueFileName string, err error) {
	// 1. generate the template file
	valueFile, err := ioutil.TempFile(os.TempDir(), "values")
	if err != nil {
		log.Errorf("Failed to create tmp file %v due to %v", valueFile.Name(), err)
		return "", err
	}

	valueFileName = valueFile.Name()
	log.Infof("Save the values file %s", valueFileName)

	// 2. dump the object into the template file
	err = toYaml(values, valueFile)
	return valueFileName, err
}

/**
* generate helm template without tiller: helm template -f values.yaml chart_name
* Exec /usr/local/bin/helm, [template -f /tmp/values313606961 --namespace default --name hj /charts/tf-horovod]
* returns generated template file: templateFileName
 */
func GenerateHelmTemplate(name string, namespace string, valueFileName string, chartName string, options ...string) (templateFileName string, err error) {
	tempName := fmt.Sprintf("%s.yaml", name)
	templateFile, err := ioutil.TempFile("", tempName)
	if err != nil {
		return templateFileName, err
	}
	templateFileName = templateFile.Name()

	binary, err := exec.LookPath(helmCmd[0])
	if err != nil {
		return templateFileName, err
	}

	// 3. check if the chart file exists
	// if strings.HasPrefix(chartName, "/") {
	if _, err = os.Stat(chartName); err != nil {
		return templateFileName, err
	}
	// }

	// 4. prepare the arguments
	args := []string{binary, "template", "-f", valueFileName,
		"--namespace", namespace,
		"--name", name,
	}
	if len(options) != 0 {
		args = append(args, options...)
	}

	args = append(args, []string{chartName, ">", templateFileName}...)

	log.Infof("Exec bash -c %v", args)

	// return syscall.Exec(cmd, args, env)
	// 5. execute the command
	log.Debugf("Generating template  %v", args)
	cmd := exec.Command("bash", "-c", strings.Join(args, " "))
	// cmd.Env = env
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", string(out))
	if err != nil {
		return templateFileName, fmt.Errorf("Failed to execute %s, %v with %v", binary, args, err)
	}

	// // 6. clean up the value file if needed
	// if log.GetLevel() != log.DebugLevel {
	// 	err = os.Remove(valueFileName)
	// 	if err != nil {
	// 		log.Warnf("Failed to delete %s due to %v", valueFileName, err)
	// 	}
	// }

	return templateFileName, nil
}

/**
* Check the chart version by given the chart directory
* helm inspect chart /charts/tf-horovod
 */

func GetChartVersion(chart string) (version string, err error) {
	binary, err := exec.LookPath(helmCmd[0])
	if err != nil {
		return "", err
	}

	// 1. check if the chart file exists, if it's it's unix path, then check if it's exist
	// if strings.HasPrefix(chart, "/") {
	if _, err = os.Stat(chart); err != nil {
		return "", err
	}
	// }

	// 2. prepare the arguments
	args := []string{binary, "inspect", "chart", chart,
		"|", "grep", "version:"}
	log.Debugf("Exec bash -c %v", args)

	cmd := exec.Command("bash", "-c", strings.Join(args, " "))
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("Failed to find version when executing %s, result is %s", args, out)
	}
	fields := strings.Split(lines[0], ":")
	if len(fields) != 2 {
		return "", fmt.Errorf("Failed to find version when executing %s, result is %s", args, out)
	}

	version = strings.TrimSpace(fields[1])
	return version, nil
}

func GetChartName(chart string) string {
	return filepath.Base(chart)
}
