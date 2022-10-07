package bootstrapnew

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	logs "github.com/openshift/odo/pkg/log"
	"github.com/redhat-developer/kam/pkg/cmd/genericclioptions"
	"github.com/redhat-developer/kam/pkg/cmd/ui"
	pipelines "github.com/redhat-developer/kam/pkg/pipelines/component"
	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const (
	DescribeRecommendedCommandName = "describe"
	appliactionNameFlags           = "application-name"
)

var (
	describeExampleC = ktemplates.Examples(`
    #
		kam describe --output <path to write GitOps resources> --application-name <name of application>
		
    %[1]s 
    `)

	describeLongDescC  = ktemplates.LongDesc(`Describe Command`)
	describeShortDescC = `Describes the details of the application `
	appFS              = ioutils.NewFilesystem()
)

func NewDescibeParameters() *DescribeParameters {
	return &DescribeParameters{
		GeneratorOptions: &pipelines.GeneratorOptions{},
	}
}

type DescribeParameters struct {
	*pipelines.GeneratorOptions
}

func nonInteractiveModeDescribe(io *DescribeParameters) error {
	mandatoryFlags := map[string]string{appliactionNameFlags: io.ApplicationName}
	if err := checkMandatoryFlagsDescribe(mandatoryFlags); err != nil {
		return err
	}
	err := ui.ValidateName(io.ApplicationName)
	if err != nil {
		return err
	}
	if io.Output == "./" || io.Output == "." {
		path, _ := os.Getwd()
		io.Output = path
	}
	exists, _ := ioutils.IsExisting(appFS, io.Output)
	if !exists {
		return fmt.Errorf("the given Path : %s  doesn't exists ", io.Output)
	}
	exists, _ = ioutils.IsExisting(appFS, filepath.Join(io.Output, io.ApplicationName))
	if !exists {
		return fmt.Errorf("the given Application: %s  doesn't exists in the Path: %s", io.ApplicationName, io.Output)
	}
	exists, _ = ioutils.IsExisting(appFS, filepath.Join(io.Output, io.ApplicationName, io.ComponentName))
	if !exists {
		return fmt.Errorf("the given Component: %s  doesn't exists in the Application: %s", io.ComponentName, io.ApplicationName)
	}
	return nil
}

func checkMandatoryFlagsDescribe(flags map[string]string) error {
	missingFlags := []string{}
	mandatoryFlags := []string{appliactionNameFlags}
	for _, flag := range mandatoryFlags {
		if flags[flag] == "" {
			missingFlags = append(missingFlags, fmt.Sprintf("%q", flag))
		}
	}
	if len(missingFlags) > 0 {
		return missingFlagErr(missingFlags)
	}
	return nil
}
func checkENv(path string) []string {
	var envList []string
	exists, _ := ioutils.IsExisting(appFS, path)
	if exists {
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			err := ui.ValidateName(f.Name())
			if err == nil {
				envList = append(envList, f.Name())
			}
		}
	} else {
		return []string{""}
	}

	return envList
}

func listFiles(path string) map[string][]string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	printList := make(map[string][]string)
	for _, f := range files {
		err = ui.ValidateName(f.Name())
		if err == nil {
			l := checkENv(filepath.Join(path, f.Name(), "overlays"))
			printList[f.Name()] = l
		}
	}
	return printList
}
func (io *DescribeParameters) Complete(name string, cmd *cobra.Command, args []string) error {

	return nonInteractiveModeDescribe(io)
}
func (io *DescribeParameters) Validate() error {
	return nil
}

func (io *DescribeParameters) Run() error {
	listComp := listFiles(filepath.Join(io.Output, io.ApplicationName, "components"))
	keys := []string{}
	for k := range listComp {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(listComp) != 0 {
		logs.Progressf("Components in application %s", io.ApplicationName)
		for f := 0; f < len(keys); f++ {
			logs.Progressf(" - %s ", keys[f])
			if listComp[keys[f]][0] != "" {
				logs.Progressf("   Environments:")
				for i := 0; i < len(listComp[keys[f]]); i++ {
					logs.Progressf("     - %s", listComp[keys[f]][i])
				}

			}
		}
	} else {
		logs.Progressf("No component is present in your application")
	}

	return nil
}

func NewCmdDescribe(name, fullName string) *cobra.Command {
	o := NewDescibeParameters()
	var descibeCmd = &cobra.Command{
		Use:     name,
		Short:   describeShortDescC,
		Long:    describeLongDescC,
		Example: fmt.Sprintf(describeExampleC, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}
	descibeCmd.Flags().StringVar(&o.Output, "output", "./", "Path to write GitOps resources")
	descibeCmd.Flags().StringVar(&o.ApplicationName, "application-name", "", "Provide a name of application")
	return descibeCmd
}
