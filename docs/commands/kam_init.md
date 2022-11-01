## kam init

Preform the Git init, branch and remote commands to intialialize the GitOps folder.

### Synopsis

Git intialize the GitOps repository.

```
kam init [flags]
```

### Examples

```
  # Init command to initialize git repository.
  kam init --application-folder <path to application> --git-repo-url <https://github.com/<your organization>/gitops.git --secret <your git access token>
  
  kam init
```

### Options

```
      --application-folder string    Provode the path to the application
      --git-repo-url string          Provide the URL for your GitOps repository e.g. https://github.com/organisation/repository.git
  -h, --help                         help for init
      --interactive                  If true, enable prompting for most options if not already specified on the command line
      --private-repo-driver string   If your Git repositories are on a custom domain, please indicate which driver to use github or gitlab
      --secret string                Used to authenticate repository clones. Access token is encrypted and stored on local file system by keyring, will be updated/reused.
```

### SEE ALSO

* [kam](kam.md)	 - kam

