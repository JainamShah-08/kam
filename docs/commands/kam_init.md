## kam init

Perform the Git init, branch and remote commands to initialize the GitOps folder.

### Synopsis

Git intialize the GitOps repository.

```
kam init [flags]
```

### Examples

```
  # Init command to initialize GitOps repository.
  kam init --application-folder <path to application> --git-repo-url <https://github.com/<your organization>/gitops.git --secret <your git access token>
  
  kam init
```

### Options

```
      --application-folder string    Provide the path to the application folder
      --git-repo-url string          Provide the URL for your GitOps repository e.g. https://github.com/organisation/repository.git
  -h, --help                         help for init
      --interactive                  If true, enable prompting for most options if not already specified on the command line
      --private-repo-driver string   If your Git repositories are on a custom domain, please indicate which driver to use github or gitlab
      --token string                 Used to authenticate repository clones. Access token is encrypted and stored on local file system by keyring, will be updated/reused.
```

### SEE ALSO

* [kam](kam.md)	 - kam
