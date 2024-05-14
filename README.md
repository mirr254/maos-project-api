Following this https://www.retgits.com/2020/01/how-to-create-a-vpc-in-aws-using-pulumi-and-golang/ for a start

Pulum over http - https://github.com/pulumi/automation-api-examples/blob/main/go/pulumi_over_http/main.go

How to run server.

### Testing locally with Microk8s
- Install and start microk8s then run `microk8s kubectl config view --raw > ~/.kube/config`
- Make changes to the code
- build the image and save image tar 
  ```sh 
    docker build -t maos-project-api:local .
    docker save maos-project-api > maos-project-api.tar
  ```
- Export the built image from the local Docker daemon and “inject” it into the MicroK8s image cache. This is done by transferring image from host to the VM managed by  
  multipass 
    ```sh
    multipass transfer maos-project-api.tar microk8s-vm:/tmp/maos-project-api.tar

    microk8s ctr image import /tmp/maos-project-api.tar
       
    ```
And it can be used in values file
```
image:
  repository: maos-project-api
  pullPolicy: Never
  # Overrides the image tag whose default is the chart appVersion.
  tag: "local"
```

Run `microk8s kubectl config view --raw > ~/.kube/config` to use kubectl normally.  

Note: Avoid tagging local images with `latest` since containerd will not cache images with that tag.
