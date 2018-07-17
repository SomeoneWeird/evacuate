package main

import "github.com/SomeoneWeird/evacuate/cmd"
import "github.com/SomeoneWeird/evacuate/plugins"
import "github.com/SomeoneWeird/evacuate/providers"

func main() {
	plugins.RegisterPlugin("kernel", plugins.Kernel{})
	plugins.RegisterPlugin("docker", plugins.Docker{})
	providers.RegisterProvider("s3", providers.S3{})

	cmd.Execute()
}
