## kam bootstrap-new

New Bootstrap Command Application Configuration

### Synopsis

New Bootstrap Command

```
kam bootstrap-new [flags]
```

### Examples

```
  # Bootstrap-New OpenShift pipelines.
  kam bootstrap-new --git-repo-url https://github.com/<your organization>/gitops.git --application-name <name of application> --component-name <name of component> --secret <your git access token> --output <path to write GitOps resources> --push-to-git=true
  
  kam bootstrap-new
```

### Options

```
      --application-name string      Provide a name for your Application
      --component-name string        Provide a Component Name within the Application
      --git-repo-url string          Provide the URL for your GitOps repository e.g. https://github.com/organisation/repository.git
  -h, --help                         help for bootstrap-new
      --interactive                  If true, enable prompting for most options if not already specified on the command line
      --namespace string             this is a name-space options (default "openshift-gitops")
      --output string                Path to write GitOps resources (default "./")
      --overwrite                    If true, it will overwrite the files
      --private-repo-driver string   If your Git repositories are on a custom domain, please indicate which driver to use github or gitlab
      --push-to-git                  Overwrites previously existing GitOps configuration (if any) on the local filesystem
      --route string                 If you specify the route flag and pass the string, that string will be in the route.yaml that is generated
      --save-token-keyring           Explicitly pass this flag to update the git-host-access-token in the keyring on your local machine
      --secret string                Used to authenticate repository clones. Access token is encrypted and stored on local file system by keyring, will be updated/reused.
      --target-port int              Provide the Target Port for your Application (default 8080)
```

### SEE ALSO

* [kam](kam.md)	 - kam

