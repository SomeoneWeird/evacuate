# evacuate

Have a misbehaving or compromised instance? Evacuate it.

## Goals

- Collect information to aid in diagnostics
- Automatically compress and upload artifacts
- Interface with running software to gracefully shut down
- Stop/Terminate the instance being run on.

## Usage

```
evacuate
```

This will detect if it's being run in a cloud environment, and if so automatically enable provider-specific plugins (eg. collecting EC2 metadata information)

By default, it will collect information and save a gzip'd file in the current working directory as `evacuate.tar.gz`. You can set up automatic uploading by passing specific configuration options, see below.

## Configuration

A configuration file can be created and passed in at runtime with `--config`.

If you want to enable automatically uploading the resulting artifact, you will need to set that up in this configuration file.

An example file may look like:

```
{
  "uploader": {
    "type": "awss3",
    "options": {
      "bucket": "my-unique-bucket"
    }
  },
  "plugins": {
    "sysInfo": {
      "enabled": false
    }
  }
}
```

## Plugins

By default, all plugins are enabled and will automatically detect if they should execute. You can override this by explicitly setting the `enabled` key for the plugin to `true`/`false` in your configuration file.

### System Information

All non-application specific information collection falls under this plugins.

Includes but not limited to:
  - Kernel Information
    - Parameters
    - Logs
  - Network
    - Interfaces
    - Firewall Rules
    - iptables dump
  - Disk Information
  - Memory Information

### Systemd

Collect status and logs of all running systemd services.

### Kubernetes

Collect information from the Kubelet running on this instance.

- Daemon Logs

### Docker

- Daemon Logs
- Running Containers

### EC2 Metadata (AWS)

If we're running as an EC2 instance, we can collect information from the metadata url located at http://169.254.169.254/latest/meta-data/ 

## Upload Providers

### AWS S3

Upload the results to an S3 bucket.

If running on EC2, instance roles are supported instead of IAM keys as long as the role has access to write to the specified bucket.

```
{
  "uploader": {
    "type": "awss3",
    "options": {
      "bucket": "my-unique-bucket"
    }
  }
}
```

#### Options

- *bucket* - **required**. The bucket name to use.
- *path* - The path inside the bucket to save the results to.
- *accessKeyId* - The access key id. Required if not using instance roles.
- *accessKeySecret* - The access key secret. Required if not using instance roles.

## Future Work

### Remote Mode

`evacuate(1)` should be able to SSH into a remote instance, download itself (for the correct architecture), and eventually sync back results to the local instance once it has run.
