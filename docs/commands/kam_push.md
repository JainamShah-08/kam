## kam push

Perform the Git add, commit and push commands for untracked changes.

### Synopsis

Pushing the application to Git repository.

```
kam push [flags]
```

### Examples

```
  # Push command to push changes to git.
  kam push --application-folder <path to application> --commit-message <Message to commit>
  
  kam push
```

### Options

```
      --application-folder string   Provode the path to the application
      --commit-message string       Provode a message to commit changes to repository
  -h, --help                        help for push
      --interactive                 If true, enable prompting for most options if not already specified on the command line
```

### SEE ALSO

* [kam](kam.md)	 - kam

