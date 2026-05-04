# Project Rules

- Never add a Co-Authored-By line to git commit messages.

## Building & Pushing tax-form-renderer to CRC Internal Registry

Registry: `default-route-openshift-image-registry.apps-crc.testing`
Containerfile: `taxform-renderer/Containerfile`
Build context: `taxform-renderer/`

```bash
# Login to the CRC internal registry
podman login -u kubeadmin -p $(oc whoami -t) default-route-openshift-image-registry.apps-crc.testing --tls-verify=false

# Build the image
podman build -t default-route-openshift-image-registry.apps-crc.testing/sandbox/tax-form-renderer:latest -f taxform-renderer/Containerfile taxform-renderer/

# Push to the CRC registry
podman push default-route-openshift-image-registry.apps-crc.testing/sandbox/tax-form-renderer:latest --tls-verify=false
```

The image will be available in-cluster as `image-registry.openshift-image-registry.svc:5000/sandbox/tax-form-renderer:latest`.
