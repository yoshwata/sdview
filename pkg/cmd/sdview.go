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
	"errors"
	"fmt"
	"os"
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
	labExample = `
	%[1]s lab pods
	%[1]s lab services
`
)

// LabOptions provides information required to update
type LabOptions struct {
	configFlags *genericclioptions.ConfigFlags

	args []string

	Namespace string

	genericclioptions.IOStreams

	output string

	sdbuildPath string
}

// NewLabOptions provides an instance of LabOptions with default values
func NewLabOptions(streams genericclioptions.IOStreams) *LabOptions {
	return &LabOptions{
		configFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}

// NewCmdLab provides a cobra command wrapping LabOptions
func NewCmdLab(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewLabOptions(streams)

	cmd := &cobra.Command{
		Use:          "kubectl lab [resources] [flags]",
		Short:        "kubectl lab",
		Example:      fmt.Sprintf(labExample, "kubectl"),
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

	var echoTimes = 1
	cmd.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")
	cmd.Flags().StringVarP(&o.output, "output", "o", "default", "hoge")
	cmd.Flags().StringVarP(&o.sdbuildPath, "sdbuildPath", "b", "default", "hoge")

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
	cases := []struct {
		want bool
		msg  string
	}{
		{
			want: len(o.args) > 0,
			msg:  "must give args",
		},
	}

	for _, c := range cases {
		if !c.want {
			return errors.New(c.msg)
		}
	}

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
	fmt.Println("faaaaaaa")
	fmt.Println(pathArg)

	templateValue := ""
	templateFormat := ""
	for format := range columnsFormats {
		format = format + "="
		if strings.HasPrefix(pathArg, format) {
			templateValue = pathArg[len(format):]
			templateFormat = format[:len(format)-1]
			break
		}
	}

	fmt.Println(templateValue)
	fmt.Println(templateFormat)

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

	columns, err := makeColumns(o.output)
	if err != nil {
		return fmt.Errorf("Failed to makeColumns")
	}
	fmt.Println("colums")
	fmt.Printf("%#v\n", columns)

	buildColumns, err := makeColumns(o.sdbuildPath)
	fmt.Println("buildColumns")
	fmt.Printf("%#v\n", buildColumns)

	k8sRes := resource.
		NewBuilder(o.configFlags).
		Unstructured().
		NamespaceParam(o.Namespace).
		DefaultNamespace().
		ResourceTypeOrNameArgs(true, o.args...).
		Latest().
		Flatten().
		Do()
	if err := k8sRes.Err(); err != nil {
		return err
	}

	many := map[string][]string{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	usertoken := viper.Get("usertoken").(string)
	sdapi := viper.Get("sdapi").(string)
	fmt.Println(usertoken)
	fmt.Println(sdapi)

	sd := screwdriver.New(usertoken, sdapi)

	if err := k8sRes.Visit(func(info *resource.Info, e error) error {
		strPod, _ := json.Marshal(info.Object)
		bytePod := []byte(strPod)

		var unMarshaledPod interface{}
		json.Unmarshal(bytePod, &unMarshaledPod)

		buildId, _ := jsonpath.Read(unMarshaledPod, "$.metadata.labels.sdbuild")
		if buildId == nil {
			return nil
		}

		var pathedPod interface{}
		for i := range columns {
			fmt.Println(columns[i].FieldSpec)
			pathedPod, _ = jsonpath.Read(unMarshaledPod, columns[i].FieldSpec)
			_ = pathedPod
			many[columns[i].Header] = append(many[columns[i].Header], pathedPod.(string))
		}

		many["buildId"] = append(many["buildId"], buildId.(string))
		// repository := "myorg/myrepo"

		sdBuild, err := sd.Build(buildId.(string))
		if err != nil {
			return nil
		}
		// sdJob := sd.Job(sdBuild.JobID)
		// sdPipeline := sd.Pipeline(sdJob.PipelineId)
		// pathedBi, _ := jsonpath.Read(sdBuild, "$.buildClusterName")
		// if pathedBi == nil {
		// 	pathedBi = ""
		// }

		for i := range buildColumns {
			fmt.Println(buildColumns[i].FieldSpec)

			pathedBi, err := jsonpath.Read(sdBuild, buildColumns[i].FieldSpec)
			if err != nil {
				pathedBi = ""
			}
			many[buildColumns[i].Header] = append(many[buildColumns[i].Header], pathedBi.(string))
		}

		// _ = sdBuild
		// many["buildClusterName"] = append(many["buildClusterName"], pathedBi.(string))

		return e
	}); err != nil {
		return err
	}

	printer := tableprinter.New(os.Stdout)
	printer.Print(many)

	return nil
}
