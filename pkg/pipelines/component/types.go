package pipelines

type GeneratorOptions struct {
	Output               string //
	ComponentName        string //
	ApplicationName      string //
	Secret               string //
	GitRepoURL           string //
	NameSpace            string //
	TargetPort           int    //
	PushToGit            bool   // If true, gitops repository is pushed to remote git repository.
	Route                string
	Overwrite            bool //
	SaveTokenKeyRing     bool
	PrivateRepoURLDriver string //
	EnvironmentName      string
	ApplicationFolder    string
}
