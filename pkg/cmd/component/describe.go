package bootstrapnew

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"

	logs "github.com/openshift/odo/pkg/log"
	"github.com/redhat-developer/kam/pkg/cmd/genericclioptions"
	"github.com/redhat-developer/kam/pkg/cmd/ui"
	pipelines "github.com/redhat-developer/kam/pkg/pipelines/component"
	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const (
	DescribeRecommendedCommandName = "describe"
	appliactionFolderFlags         = "application-folder"
)

var (
	describeExampleC = ktemplates.Examples(`
    # Describe the application
		kam describe --application-folder <path to application>
		
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
	mandatoryFlags := map[string]string{appliactionFolderFlags: io.ApplicationFolder}
	if err := CheckMandatoryFlags(mandatoryFlags); err != nil {
		return err
	}
	exists, _ := ioutils.IsExisting(appFS, io.ApplicationFolder)
	if !exists {
		return fmt.Errorf("the given Path : %s  doesn't exists ", io.ApplicationFolder)
	}
	exists, _ = ioutils.IsExisting(appFS, filepath.Join(io.ApplicationFolder, "components"))
	if !exists {
		return fmt.Errorf("the given Path : %s is not a correct path for an application ", io.ApplicationFolder)
	}
	return nil
}

// func DummyFalgCheck()error{

// }

func checkEnv(currentFileSystem afero.Afero, path string) []string {
	var envList []string
	exists, _ := ioutils.IsExisting(currentFileSystem, path)
	if exists {
		files, _ := currentFileSystem.ReadDir(path)
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

func ListFiles(appFSvar afero.Afero, path string) map[string][]string {
	files, err := appFSvar.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	printList := make(map[string][]string)
	for _, f := range files {
		err = ui.ValidateName(f.Name())
		if err == nil {
			l := checkEnv(appFS, filepath.Join(path, f.Name(), "overlays"))
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
	listComp := ListFiles(appFS, filepath.Join(io.ApplicationFolder, "components"))
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
	descibeCmd.Flags().StringVar(&o.ApplicationFolder, "application-folder", "", "Provide the path to the application")
	return descibeCmd
}
