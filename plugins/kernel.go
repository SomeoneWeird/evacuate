package plugins

// Kernel collects information about the host kernel
type Kernel struct{}

// ShouldRun ensures we're on a platform we support
func (p Kernel) ShouldRun(ctx PluginContext) {
	ctx.Finish <- true
}

// Run executes this plugin
func (p Kernel) Run(ctx PluginContext) {
	ctx.Finish <- true
}
