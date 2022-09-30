## kam component delete

Delete an existing Component

### Synopsis

Delete an existing Component from the Application

```
kam component delete [flags]
```

### Examples

```
  # Delete an existing Component from the Application
  # Example: kam component delete --output <path to Application folder> --application-name <Application name to delete component> --component-name <name of the component to delete>
  
  kam component delete
```

### Options

```
      --application-name string   Name of the Application to delete a Component
      --component-name string     Name of the component to delete
  -h, --help                      help for delete
      --interactive               If true, enable prompting for most options if not already specified on the command line
      --output string             Folder path to the Application to delete the Component (default "./")
```

### SEE ALSO

* [kam component](kam_component.md)	 - Manage component in application

