package bootstrapnew

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/openshift/odo/pkg/log"
	gitops "github.com/redhat-developer/gitops-generator/pkg"
	"github.com/redhat-developer/kam/pkg/cmd/component/cmd/ui"
	"github.com/redhat-developer/kam/pkg/cmd/genericclioptions"
	"github.com/redhat-developer/kam/pkg/cmd/utility"
	"github.com/redhat-developer/kam/pkg/pipelines/accesstoken"
	pipelines "github.com/redhat-developer/kam/pkg/pipelines/component"
	"github.com/redhat-developer/kam/pkg/pipelines/ioutils"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

const (
	InitRecommendedCommandName = "init"
	defaultRepoDescription     = "Bootstrapped GitOps Repository based on Components"
	applicationFolderFlag      = "application-folder"
)

var (
	initExample = ktemplates.Examples(`
   # Init command to initialize GitOps repository.
       kam init --application-folder <path to application> --git-repo-url <https://github.com/<your organization>/gitops.git --secret <your git access token>
      
   %[1]s
   `)

	initLongDesc  = ktemplates.LongDesc(`Git intialize the GitOps repository.`)
	initShortDesc = `Preform the Git init, branch and remote commands to initialize the GitOps folder.`
)

func NewInitParameters() *InitParameters {
	return &InitParameters{
		GeneratorOptions: &pipelines.GeneratorOptions{},
	}
}

type InitParameters struct {
	*pipelines.GeneratorOptions
	Interactive bool
}

// checks the application folder constraints for non interactive mode
// if the application folder path is present then it check wether it is valid for application
// if the folder path is not present check for valid application name to create it

func CheckApplicationPath(app afero.Afero, path string) error {
	exists, _ := ioutils.IsExisting(app, path)
	if !exists {
		appName := strings.Split(path, "/")
		err := ui.ValidateName(appName[len(appName)-1])
		if err != nil {
			return fmt.Errorf("failed to create the directory %s with application name %s : %v", path, appName[len(appName)-1], err)
		}
		app.MkdirAll(path, 0755)
	}
	subFiles := ListFiles(app, path)
	if len(subFiles) != 0 {
		exists, _ = ioutils.IsExisting(app, filepath.Join(path, "components"))
		if !exists {
			return fmt.Errorf("the given path %s is not the correct path for application", path)
		}
	}
	return nil
}
func nonInteractiveModeInit(io *InitParameters) error {
	passedFlags := map[string]string{appliactionFolderFlags: io.ApplicationFolder}
	if err := checkMandatoryFlagsDescribe(passedFlags); err != nil {
		return err
	}

	if err := CheckApplicationPath(appFS, io.ApplicationFolder); err != nil {
		return err
	}
	if io.PrivateRepoURLDriver != "" {
		if !supportedDrivers.supported(io.PrivateRepoURLDriver) {
			return fmt.Errorf("invalid driver type: %q", io.PrivateRepoURLDriver)
		}
	}
	if checkGit(appFS, io.ApplicationFolder) {
		return fmt.Errorf("git repository has already been initiated")
	} else {
		if io.GitRepoURL == "" {
			return fmt.Errorf("git repository cannot be initiated without git-repo-url")
		} else {
			_, err := url.Parse(io.GitRepoURL)
			if err != nil {
				return fmt.Errorf("failed to parse GitOps repo URL %q: %v", io.GitRepoURL, err)
			}
			err = setAccessTokenInit(io)
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
			return gitInitializeCheck(io)
		}
	}
}

func interactiveModeInit(io *InitParameters) error {
	log.Progressf("\nStarting interactive prompt\n")
	var url *url.URL
	var err error
	if io.ApplicationFolder != "" {
		//check wether the path exists and if not check the application-name
		exists, _ := ioutils.IsExisting(appFS, io.ApplicationFolder)
		if !exists {
			appName := strings.Split(io.ApplicationFolder, "/")
			err := ui.ValidateName(appName[len(appName)-1])
			if err != nil {
				log.Progressf("failed to create the directory %s with application name %s : %v", io.ApplicationFolder, appName[len(appName)-1], err)
				io.ApplicationFolder = ui.ApplicationOutputPath()
			}
		} else {
			listDir := ListFiles(appFS, io.ApplicationFolder)
			if len(listDir) != 0 {
				exists, _ = ioutils.IsExisting(appFS, filepath.Join(io.ApplicationFolder, "components"))
				if !exists {
					log.Progressf("the given Path : %s is not a correct path for an application ", io.ApplicationFolder)
					io.ApplicationFolder = ui.ApplicationOutputPath()
				}
			}
		}

	}
	if io.ApplicationFolder == "" {
		io.ApplicationFolder = ui.ApplicationOutputPath()
	}
	// ask for confirmation before creating the directory
	exists, _ := ioutils.IsExisting(appFS, io.ApplicationFolder)
	if !exists {
		check := ui.AskConfirmation(io.ApplicationFolder)
		if check {
			appFS.MkdirAll(io.ApplicationFolder, 0755)
		} else {
			return fmt.Errorf("interactive mode has been terminated")
		}
	}

	if !checkGit(appFS, io.ApplicationFolder) {
		if io.GitRepoURL == "" {
			io.GitRepoURL = ui.EnterGitRepoURL()
		}
		url, err = url.Parse(io.GitRepoURL)
		for err != nil {
			log.Progressf("failed to parse GitOps repo URL %q: %v", io.GitRepoURL, err)
			io.GitRepoURL = ui.EnterGitRepoURL()
			url, err = url.Parse(io.GitRepoURL)
		}
		io.GitRepoURL = utility.AddGitSuffixIfNecessary(io.GitRepoURL)
		if !isKnownDriverURL(io.GitRepoURL) {
			io.PrivateRepoURLDriver = ui.SelectPrivateRepoDriver()
			host, err := accesstoken.HostFromURL(io.GitRepoURL)
			if err != nil {
				return fmt.Errorf("failed to parse the gitops url: %v", err)
			}
			identifier := factory.NewDriverIdentifier(factory.Mapping(host, io.PrivateRepoURLDriver))
			factory.DefaultIdentifier = identifier
		}
		secret, err := accesstoken.GetAccessToken(io.GitRepoURL)
		if err != nil && err != keyring.ErrNotFound {
			return err
		}
		if secret == "" { // We must prompt for the token
			if io.Secret == "" {
				io.Secret = ui.EnterGitSecret(io.GitRepoURL)
				io.SaveTokenKeyRing = ui.UseKeyringRingSvc()
				if io.SaveTokenKeyRing {
					setAccessTokenInit(io)
				}
			}
		}
		io.Secret = secret
	} else {
		return fmt.Errorf("git repo already been initiated")
	}
	return gitInitializeCheck(io)
}

// Checks the requirement for .git folder
func gitInitializeCheck(io *InitParameters) error {
	u, _ := url.Parse(io.GitRepoURL)
	parts := strings.Split(u.Path, "/")
	org := parts[1]
	repoName := strings.TrimSuffix(strings.Join(parts[2:], "/"), ".git")
	u.User = url.UserPassword("", io.Secret)
	client, err := factory.FromRepoURL(u.String())
	if err != nil {
		return fmt.Errorf("failed to create a client to access %q: %v", io.GitRepoURL, err)
	}
	ctx := context.Background()
	currentUser, _, err := client.Users.Find(ctx)
	if err != nil {
		return fmt.Errorf("failed to get the user with their auth token: %v", err)
	}
	if currentUser.Login == org {
		org = ""
	}
	ri := &scm.RepositoryInput{
		Private:     true,
		Description: defaultRepoDescription,
		Namespace:   org,
		Name:        repoName,
	}
	_, _, err = client.Repositories.Create(context.Background(), ri)
	if err != nil {
		repo := fmt.Sprintf("%s/%s", org, repoName)
		if org == "" {
			repo = fmt.Sprintf("%s/%s", currentUser.Login, repoName)
		}
		if _, resp, err := client.Repositories.Find(context.Background(), repo); err == nil && resp.Status == 200 {
			return fmt.Errorf("failed to create repository, repo already exists")
		}
		return fmt.Errorf("failed to create repository %q in namespace %q: %v", repoName, org, err)
	}
	return gitInitialize(io.ApplicationFolder, io.GitRepoURL)
}

// Generates the .git folder
// Executes git init, git branch -m main and git remote add commands
func gitInitialize(path string, repo string) error {
	e := gitops.NewCmdExecutor()
	if out, err := e.Execute(path, "git", "init", "."); err != nil {
		return fmt.Errorf("failed to initialize git repository in %q %q: %s", path, string(out), err)
	}
	if out, err := e.Execute(path, "git", "branch", "-m", "main"); err != nil {
		return fmt.Errorf("failed to switch to branch %q in repository in %q %q: %s", "main", path, string(out), err)
	}
	if out, err := e.Execute(path, "git", "remote", "add", "origin", repo); err != nil {
		return fmt.Errorf("failed to add files for component , to remote 'origin'  to repository in %q %q: %s", path, string(out), err)
	}
	nextStepsInit()
	return nil
}

// check the pre existence of .git folder for application folder
func checkGit(appFS afero.Afero, path string) bool {
	exists, _ := ioutils.IsExisting(appFS, filepath.Join(path, ".git"))
	if exists {
		return exists
	}
	return false
}

func (io *InitParameters) Complete(name string, cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() == 0 || io.Interactive {
		return interactiveModeInit(io)
	}
	return nonInteractiveModeInit(io)
}

func (io *InitParameters) Validate() error {
	return nil
}

func (io *InitParameters) Run() error {
	return nil
}
func Init(name, fullName string) *cobra.Command {
	o := NewInitParameters()
	initCmd := &cobra.Command{
		Use:     "init",
		Short:   initShortDesc,
		Long:    initLongDesc,
		Example: fmt.Sprintf(initExample, fullName),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(o, cmd, args)
		},
	}
	initCmd.Flags().StringVar(&o.ApplicationFolder, "application-folder", "", "Provide the path to the application folder")
	initCmd.Flags().BoolVar(&o.Interactive, "interactive", false, "If true, enable prompting for most options if not already specified on the command line")
	initCmd.Flags().StringVar(&o.Secret, "secret", "", "Used to authenticate repository clones. Access token is encrypted and stored on local file system by keyring, will be updated/reused.")
	initCmd.Flags().StringVar(&o.GitRepoURL, "git-repo-url", "", "Provide the URL for your GitOps repository e.g. https://github.com/organisation/repository.git")
	initCmd.Flags().StringVar(&o.PrivateRepoURLDriver, "private-repo-driver", "", "If your Git repositories are on a custom domain, please indicate which driver to use github or gitlab")
	return initCmd
}
func setAccessTokenInit(io *InitParameters) error {
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
func nextStepsInit() {
	log.Success("Successfully created the git repository\n\n",
		"\n",
	)
}
