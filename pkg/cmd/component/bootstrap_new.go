package bootstrapnew

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/openshift/odo/pkg/log"
	"github.com/redhat-developer/kam/pkg/cmd/component/cmd/ui"
	"github.com/redhat-developer/kam/pkg/cmd/genericclioptions"
	"github.com/redhat-developer/kam/pkg/cmd/utility"
	"github.com/redhat-developer/kam/pkg/pipelines/accesstoken"
	pipelines "github.com/redhat-developer/kam/pkg/pipelines/component"
	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

//testing changes
const (
	BootstrapRecommendedCommandName = "bootstrap-new"

	componentNameFlag     = "component-name"
	applicationNameFlag   = "application-name"
	gitRepoURLFlag        = "git-repo-url"
	secretFlag            = "secret"
	applicationFolderFlag = "application-folder"
	commitMessageFlag     = "commit-message"
)

var (
	bootstrapExampleC = ktemplates.Examples(`
    # Bootstrap-New OpenShift pipelines.
		kam bootstrap-new --git-repo-url https://github.com/<your organization>/gitops.git --application-name <name of application> --component-name <name of component> --secret <your git access token> --output <path to write GitOps resources> --push-to-git=true
		
    %[1]s 
    `)

	bootstrapLongDescC  = ktemplates.LongDesc(`New Bootstrap Command`)
	bootstrapShortDescC = `New Bootstrap Command Application Configuration`
)

// BootstrapNewParameters encapsulates the parameters for the kam pipelines init command.
type BootstrapNewParameters struct {
	*pipelines.GeneratorOptions
	Interactive bool
}

type drivers []string

var (
	supportedDrivers = drivers{
		"github",
		"gitlab",
	}
)

func (d drivers) supported(s string) bool {
	for _, v := range d {
		if s == v {
			return true
		}
	}
	return false
}

// NewBootstrapNewParameters bootsraps a Bootstrap Parameters instance.
func NewBootstrapNewParameters() *BootstrapNewParameters {
	return &BootstrapNewParameters{
		GeneratorOptions: &pipelines.GeneratorOptions{},
	}
}

// Complete completes BootstrapNewParameters after they've been created.
// If the prefix provided doesn't have a "-" then one is added, this makes the
// generated environment names nicer to read.
func (io *BootstrapNewParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	_, err := utility.NewClient()
	if err != nil {
		return err
	}

	if io.PrivateRepoURLDriver != "" {
		host, err := accesstoken.HostFromURL(io.GitRepoURL)
		if err != nil {
			return err
		}
		identifier := factory.NewDriverIdentifier(factory.Mapping(host, io.PrivateRepoURLDriver))
		factory.DefaultIdentifier = identifier
	}

	if cmd.Flags().NFlag() == 0 || io.Interactive {
		return initiateInteractiveModeForBootstrapNewCommand(io, cmd)
	}
	addGitURLSuffixIfNecessary(io)
	return nonInteractiveModeBootstrapNew(io)
}
func addGitURLSuffixIfNecessary(io *BootstrapNewParameters) {
	io.GitRepoURL = utility.AddGitSuffixIfNecessary(io.GitRepoURL)
}

// nonInteractiveMode gets triggered if a flag is passed, checks for mandatory flags.
func nonInteractiveModeBootstrapNew(io *BootstrapNewParameters) error {
	mandatoryFlags := map[string]string{componentNameFlag: io.ComponentName, applicationNameFlag: io.ApplicationName, gitRepoURLFlag: io.GitRepoURL, secretFlag: io.Secret}
	if err := CheckMandatoryFlags(mandatoryFlags); err != nil {
		return err
	}
	err := ui.ValidateTargetPort(io.TargetPort)
	if err != nil {
		return fmt.Errorf("%v Target Port is not valid", io.TargetPort)
	}

	err = ui.ValidateName(io.ComponentName)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	err = ui.ValidateName(io.ApplicationName)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if io.PrivateRepoURLDriver != "" {
		if !supportedDrivers.supported(io.PrivateRepoURLDriver) {
			return fmt.Errorf("invalid driver type: %q", io.PrivateRepoURLDriver)
		}
	}
	err = setAccessToken(io)
	if err != nil {
		return err
	}

	return nil
}

// Checking the mandatory flags & reusing missingFlagErr in Bootstrap.go
func CheckMandatoryFlags(flags map[string]string) error {
	missingFlags := []string{}
	mandatoryFlags := []string{}
	for k := range flags {
		mandatoryFlags = append(mandatoryFlags, k)
	}
	sort.Strings(mandatoryFlags)
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

func missingFlagErr(flags []string) error {
	return fmt.Errorf("required flag(s) %s not set", strings.Join(flags, ", "))
}

//Interactive mode for Bootstrap-mew Command
func initiateInteractiveModeForBootstrapNewCommand(io *BootstrapNewParameters, cmd *cobra.Command) error {
	log.Progressf("\nStarting interactive prompt\n")
	//Checks for mandatory flags
	promp := !ui.UseDefaultValues()
	if io.ApplicationName != "" {
		err := ui.ValidateName(io.ApplicationName)
		if err != nil {
			log.Progressf("%v Application Name is not valid %v", io.ApplicationName, err)
			io.ApplicationName = ui.AddApplicationName()
		}
	}
	if io.ApplicationName == "" {
		io.ApplicationName = ui.AddApplicationName()
	}
	if io.ComponentName != "" {
		err := ui.ValidateName(io.ComponentName)
		if err != nil {
			log.Progressf("%v Component Name is not valid %v", io.ComponentName, err)
			io.ComponentName = ui.AddComponentName()
		}
	}
	if io.ComponentName == "" {
		io.ComponentName = ui.AddComponentName()
	}
	if io.GitRepoURL == "" {
		io.GitRepoURL = ui.EnterGitRepoURL()
	}
	io.GitRepoURL = utility.AddGitSuffixIfNecessary(io.GitRepoURL)
	if !isKnownDriverURL(io.GitRepoURL) {
		io.PrivateRepoURLDriver = ui.SelectPrivateRepoDriver()
		host, err := accesstoken.HostFromURL(io.GitRepoURL)
		if err != nil {
			return fmt.Errorf("failed to parse the gitops url: %w", err)
		}
		identifier := factory.NewDriverIdentifier(factory.Mapping(host, io.PrivateRepoURLDriver))
		factory.DefaultIdentifier = identifier
	}
	// We are checking if any existing token is present.
	//If not we ask the uer to pass the token.
	//EnterGitSecret is just validating length for now.
	secret, err := accesstoken.GetAccessToken(io.GitRepoURL)
	if err != nil && err != keyring.ErrNotFound {
		return err
	}
	if secret == "" { // We must prompt for the token
		if io.Secret == "" {
			io.Secret = ui.EnterGitSecret(io.GitRepoURL)
		}
		if !cmd.Flag("save-token-keyring").Changed {
			io.SaveTokenKeyRing = ui.UseKeyringRingSvc()
		}
		setAccessToken(io)
	} else {
		io.Secret = secret
	}

	//Optional flags
	if !cmd.Flag("target-port").Changed && promp {
		io.TargetPort = ui.AddTargetPort()
	}
	if !cmd.Flag("route").Changed && promp {
		io.Route = ui.SelectRoute()
	}
	if !cmd.Flag("push-to-git").Changed && promp {
		io.PushToGit = ui.SelectOptionPushToGit()
	}

	outputPathOverridden := cmd.Flag("output").Changed

	appFs := ioutils.NewFilesystem()
	io.Output, io.Overwrite = ui.VerifyOutput(appFs, io.Output, io.Overwrite, io.ApplicationName, outputPathOverridden, promp)
	if !io.Overwrite {
		if ui.PathExists(appFs, filepath.Join(io.Output, io.ApplicationName)) {
			return fmt.Errorf("the secrets folder located as a sibling of the output folder %s already exists. Delete or rename the secrets folder and try again", io.Output)
		}
		if io.PushToGit && ui.PathExists(appFs, filepath.Join(io.Output, ".git")) {
			return fmt.Errorf("the .git folder in output path %s already exists. Delete or rename the .git folder and try again", io.Output)
		}
	}

	return nil
}

func setAccessToken(io *BootstrapNewParameters) error {
	if io.SaveTokenKeyRing {
		err := accesstoken.SetAccessToken(io.GitRepoURL, io.Secret)
		if err != nil {
			return err
		}
	}
	if io.Secret == "" {
		secret, err := accesstoken.GetAccessToken(io.GitRepoURL)
		if err != nil {
			return fmt.Errorf("unable to use access-token from keyring/env-var: %v, please pass a valid token to --git-host-access-token", err)
		}
		io.Secret = secret
	}
	return nil
}

// Validate validates the parameters of the CompomemtParameters.
func (io *BootstrapNewParameters) Validate() error {

	gr, err := url.Parse(io.GitRepoURL)
	if err != nil {
		return fmt.Errorf("failed to parse url %s: %w", io.GitRepoURL, err)
	}

	if len(utility.RemoveEmptyStrings(strings.Split(gr.Path, "/"))) != 2 {
		return fmt.Errorf("repo must be org/repo: %s", strings.Trim(gr.Path, ".git"))
	}

	if io.PrivateRepoURLDriver != "" {
		if !supportedDrivers.supported(io.PrivateRepoURLDriver) {
			return fmt.Errorf("invalid driver type: %q", io.PrivateRepoURLDriver)
		}
	}

	if io.SaveTokenKeyRing && io.Secret == "" {
		return errors.New("--secret is required if --save-token-keyring is enabled")
	}
	return nil
}

// Run runs the project Component command.
func (io *BootstrapNewParameters) Run() error {
	log.Progressf("\nCompleting Bootstrap process\n")
	appFs := ioutils.NewFilesystem()
	err := pipelines.BootstrapNew(io.GeneratorOptions, appFs)
	if err != nil {
		return err
	}
	if err == nil && io.PushToGit {
		log.Successf("Created repository")
	}
	nextSteps()
	return nil
}

func NewCmdBootstrapNew(name, fullName string) *cobra.Command {
	o := NewBootstrapNewParameters()
	var bootstrapCmd = &cobra.Command{
		Use:     name,
		Short:   bootstrapShortDescC,
		Long:    bootstrapLongDescC,
		Example: fmt.Sprintf(bootstrapExampleC, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}

	bootstrapCmd.Flags().StringVar(&o.Output, "output", ".", "Path to write GitOps resources")
	bootstrapCmd.Flags().StringVar(&o.ComponentName, "component-name", "", "Provide a Component Name within the Application")
	bootstrapCmd.Flags().StringVar(&o.ApplicationName, "application-name", "", "Provide a name for your Application")
	bootstrapCmd.Flags().StringVar(&o.Secret, "secret", "", "Used to authenticate repository clones. Access token is encrypted and stored on local file system by keyring, will be updated/reused.")
	bootstrapCmd.Flags().StringVar(&o.GitRepoURL, "git-repo-url", "", "Provide the URL for your GitOps repository e.g. https://github.com/organisation/repository.git")
	bootstrapCmd.Flags().StringVar(&o.NameSpace, "namespace", "openshift-gitops", "this is a name-space options")
	bootstrapCmd.Flags().IntVar(&o.TargetPort, "target-port", 8080, "Provide the Target Port for your application's component")
	bootstrapCmd.Flags().BoolVar(&o.PushToGit, "push-to-git", false, "Overwrites previously existing GitOps configuration (if any) on the local filesystem")
	bootstrapCmd.Flags().StringVar(&o.Route, "route", "", "Provide the route to expose the component with. If provided, it will be referenced in the generated route.yaml")
	bootstrapCmd.Flags().BoolVar(&o.Interactive, "interactive", false, "If true, enable prompting for most options if not already specified on the command line")
	bootstrapCmd.Flags().BoolVar(&o.Overwrite, "overwrite", false, "If true, it will overwrite the files")
	bootstrapCmd.Flags().BoolVar(&o.SaveTokenKeyRing, "save-token-keyring", false, "Explicitly pass this flag to update the git-host-access-token in the keyring on your local machine")
	bootstrapCmd.Flags().StringVar(&o.PrivateRepoURLDriver, "private-repo-driver", "", "If your Git repositories are on a custom domain, please indicate which driver to use github or gitlab")
	return bootstrapCmd
}

func isKnownDriverURL(repoURL string) bool {
	host, err := accesstoken.HostFromURL(repoURL)
	if err != nil {
		return false
	}
	_, err = factory.DefaultIdentifier.Identify(host)
	return err == nil
}

func nextSteps() {
	log.Success("New Bootstrapped OpenShift resources successfully\n\n",
		"\n",
	)
}
