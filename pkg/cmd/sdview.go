/*
Copyright 2018 The Kubernetes Authors.

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

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/resource"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/lensesio/tableprinter"
	"github.com/spf13/viper"
	"github.com/yalp/jsonpath"

	"github.com/yoshwata/sdview/pkg/screwdriver"
)

var (
	sdviewExample = `
	%[1]s -o="custom-columns=NAME:$.metadata.name,IMAGE:$.spec.containers[0].image" -b="custom-columns=builcClusterName:$.buildClusterName" -j="custom-columns=jobname:$.name" -e="custom-columns=causeMessage:$.causeMessage" -p="custom-columns=REPO:$.scmRepo.name"
`
)

// LabOptions provides information required to update
type LabOptions struct {
	configFlags *genericclioptions.ConfigFlags

	args []string

	Namespace string

	genericclioptions.IOStreams

	output string

	sdBuildPath string

	sdEventPath string

	sdJobPath string

	sdPipelinePath string

	maxLines int
}

// NewLabOptions provides an instance of LabOptions with default values
func NewSdViewOptions(streams genericclioptions.IOStreams) *LabOptions {
	return &LabOptions{
		configFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}

// NewCmdLab provides a cobra command wrapping LabOptions
func NewCmdSdView(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewSdViewOptions(streams)

	cmd := &cobra.Command{
		Use:          "sdview [flags]",
		Short:        "sdview",
		Example:      fmt.Sprintf(sdviewExample, "sdview"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	o.configFlags.AddFlags(cmd.Flags())

	cmd.Flags().StringVarP(&o.output, "output", "o", "", "Path of kubernetes pods response")
	cmd.Flags().StringVarP(&o.sdBuildPath, "sdbuildPath", "b", "", "Path of sd's /builds response")
	cmd.Flags().StringVarP(&o.sdJobPath, "sdJobPath", "j", "", "Path of sd's /jobs response")
	cmd.Flags().StringVarP(&o.sdEventPath, "sdEventPath", "e", "", "Path of sd's /events response")
	cmd.Flags().StringVarP(&o.sdPipelinePath, "sdPipelinePath", "p", "", "Path of sd's /pipelines response")
	cmd.Flags().IntVarP(&o.maxLines, "maxLines", "l", 0, "Max lines of table.")

	return cmd
}

// Complete sets all information required for updating the current context
func (o *LabOptions) Complete(cmd *cobra.Command, args []string) (err error) {
	o.args = args
	o.Namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	return err
}

// Validate ensures that all required arguments and flag values are provided
func (o *LabOptions) Validate() error {
	// nothing to do
	// cases := []struct {
	// 	want bool
	// 	msg  string
	// }{
	// 	{
	// 		want: len(o.args) > 0,
	// 		msg:  "must give args",
	// 	},
	// }

	// for _, c := range cases {
	// 	if !c.want {
	// 		return errors.New(c.msg)
	// 	}
	// }

	return nil
}

var columnsFormats = map[string]bool{
	"custom-columns-file": true,
	"custom-columns":      true,
}

type Column struct {
	// The header to print above the column, general style is ALL_CAPS
	Header string
	// The pointer to the field in the object to print in JSONPath form
	// e.g. {.ObjectMeta.Name}, see pkg/util/jsonpath for more details.
	FieldSpec string
}

func makeColumns(pathArg string) ([]Column, error) {
	templateValue := ""
	for format := range columnsFormats {
		format = format + "="
		if strings.HasPrefix(pathArg, format) {
			templateValue = pathArg[len(format):]
			// templateFormat = format[:len(format)-1]
			break
		}
	}

	parts := strings.Split(templateValue, ",")
	columns := make([]Column, len(parts))

	for ix := range parts {
		colSpec := strings.SplitN(parts[ix], ":", 2)
		if len(colSpec) != 2 {
			return nil, fmt.Errorf("unexpected custom-columns spec: %s, expected <header>:<json-path-expr>", parts[ix])
		}

		columns[ix] = Column{Header: colSpec[0], FieldSpec: colSpec[1]}
	}

	return columns, nil
}

// Run lists all available namespaces on a user's KUBECONFIG or updates the
// current context based on a provided namespace.
func (o *LabOptions) Run() error {

	var err error
	var kubeColumns []Column
	if o.output != "" {
		// kube columns
		kubeColumns, err = makeColumns(o.output)
		if err != nil {
			return fmt.Errorf("Failed to makeColumns: %s", err)
		}
	}

	// build columns
	buildColumns, err := makeColumns(o.sdBuildPath)

	// event columns
	eventColumns, err := makeColumns(o.sdEventPath)

	// job columns
	jobColumns, err := makeColumns(o.sdJobPath)

	// pipeline columns
	pipelineColumns, err := makeColumns(o.sdPipelinePath)

	k8sRes := resource.
		NewBuilder(o.configFlags).
		Unstructured().
		NamespaceParam(o.Namespace).
		DefaultNamespace().
		ResourceTypeOrNameArgs(true, "pods").
		Latest().
		Flatten().
		Do()
	if err := k8sRes.Err(); err != nil {
		return err
	}

	many := map[string][]string{}

	viper.SetConfigName("sdview_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/")
	viper.ReadInConfig()
	usertoken := viper.Get("usertoken").(string)
	sdapi := viper.Get("sdapi").(string)

	sd := screwdriver.New(usertoken, sdapi)

	lineCount := 0

	if err := k8sRes.Visit(func(info *resource.Info, e error) error {

		if o.maxLines != 0 && lineCount >= o.maxLines {
			return nil
		}

		lineCount += 1

		// Create table
		// Get buildId to show other elements like events, jobs, piepelines
		strPod, _ := json.Marshal(info.Object)
		bytePod := []byte(strPod)

		var unMarshaledPod interface{}
		json.Unmarshal(bytePod, &unMarshaledPod)

		buildId, _ := jsonpath.Read(unMarshaledPod, "$.metadata.labels.sdbuild")
		if buildId == nil {
			return nil
		}
		// kube
		many = appendSDInfo(unMarshaledPod, many, kubeColumns)

		many["buildId"] = append(many["buildId"], buildId.(string))
		sdBuild := sd.Build(buildId.(string))

		// event
		eventIdf, _ := jsonpath.Read(sdBuild, "$.eventId")
		if eventIdf == nil {
			return nil
		}

		eventId := strconv.FormatFloat(eventIdf.(float64), 'f', -1, 64)
		sdEvent := sd.Events(eventId)

		many = appendSDInfo(sdEvent, many, eventColumns)

		// build
		many = appendSDInfo(sdBuild, many, buildColumns)

		// job
		jobIdf, err := jsonpath.Read(sdBuild, "$.jobId")
		if err != nil {
			return nil
		}
		jobId := strconv.FormatFloat(jobIdf.(float64), 'f', -1, 64)
		sdJob := sd.Job(jobId)

		many = appendSDInfo(sdJob, many, jobColumns)

		// pipeline
		pipelineIdf, err := jsonpath.Read(sdJob, "$.pipelineId")
		if err != nil {
			return nil
		}
		pipelineId := strconv.FormatFloat(pipelineIdf.(float64), 'f', -1, 64)
		sdPipeline := sd.Pipeline(pipelineId)
		many = appendSDInfo(sdPipeline, many, pipelineColumns)

		return e
	}); err != nil {
		return err
	}

	printer := tableprinter.New(os.Stdout)
	printer.Print(many)

	return nil
}

func appendSDInfo(
	targetSource interface{},
	targetDest map[string][]string,
	columns []Column) map[string][]string {

	for i := range columns {
		pathedPi, err := jsonpath.Read(targetSource, columns[i].FieldSpec)
		if err != nil {
			pathedPi = ""
		}

		var hoge interface{}
		switch pathedPi := pathedPi.(type) {
		case float64:
			hoge = strconv.FormatFloat(pathedPi, 'f', -1, 64)
		case string:
			hoge = pathedPi
		}
		targetDest[columns[i].Header] = append(targetDest[columns[i].Header], hoge.(string))
	}

	return targetDest
}
