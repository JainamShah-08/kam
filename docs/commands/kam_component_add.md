## kam component add

Add a new Component

### Synopsis

Add a new Component to the Application

```
kam component add [flags]
```

### Examples

```
  # Add a new Component to Application
  # Example: kam component add --output <path to Application folder> --application-name <Application name to add component> --component-name new-component
  
  kam component add
```

### Options

```
      --application-name string   Name of the Application to add a Component
      --component-name string     Name of the component
  -h, --help                      help for add
      --interactive               If true, enable prompting for most options if not already specified on the command line
      --output string             Folder path to the Application to add the Component (default "./")
      --route string              Provide the route to expose the component with. If provided, it will be referenced in the generated route.yaml
      --target-port int           Provide the Target Port for your Application (default 8080)
```

### SEE ALSO

* [kam component](kam_component.md)	 - Manage component in application

