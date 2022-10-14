## kam env add

Add a new environment

### Synopsis

Add a new environment to the application's component in the GitOps repository

```
kam env add [flags]
```

### Examples

```
  # Add a new environment to the application's component in the GitOps repository
  # Example: kam env add --output <path to Application folder> --application-name <Application name> --component-name <component name> --env-name <environment name>
  
  kam env add
```

### Options

```
      --application-name string   Name of the Application to add the environment
      --component-name string     Name of the component to add the environment
      --env-name string           Name of the environment/namespace
  -h, --help                      help for add
      --interactive               If true, enable prompting for most options if not already specified on the command line
      --output string             Folder path to the Application to add a new environment (default ".")
```

### SEE ALSO

* [kam env](kam_env.md)	 - Manage an environment in GitOps

