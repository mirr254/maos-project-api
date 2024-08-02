
package aws_helpers

import (
	"github.com/pulumi/pulumi-eks/sdk/v2/go/eks"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/pulumi/pulumi-awsx/sdk/v2/go/awsx/ec2"
)

// CreateEKSResource: function is responsible for creating a new EKS cluster in the project stack.
func CreateEKSResource(ctx *pulumi.Context, project_name string, vpc *ec2.Vpc) (*eks.Cluster, error) {
	cfg := config.New(ctx, "")

	minClusterSize, err := cfg.TryInt("minClusterSize")
	if err != nil {
		minClusterSize = 3
	}
	maxClusterSize, err := cfg.TryInt("maxClusterSize")
	if err != nil {
		maxClusterSize = 5
	}
	desiredClusterSize, err := cfg.TryInt("desiredClusterSize")
	if err != nil {
		desiredClusterSize = 3
	}
	eksNodeInstanceType, err := cfg.Try("eksNodeInstanceType")
	if err != nil {
		eksNodeInstanceType = "t2.medium"
	}

	// Create an EKS cluster.
	cluster, err := eks.NewCluster(ctx, project_name, &eks.ClusterArgs{
		InstanceType: pulumi.String(eksNodeInstanceType),
		VpcId: 	  vpc.VpcId,
		PublicSubnetIds: vpc.PublicSubnetIds,
		PrivateSubnetIds: vpc.PrivateSubnetIds,
		MinSize: pulumi.Int(minClusterSize),
		MaxSize: pulumi.Int(maxClusterSize),
		DesiredCapacity: pulumi.Int(desiredClusterSize),
		NodeAssociatePublicIpAddress: pulumi.BoolRef(false),

		// Change these values for a private cluster (VPN access required)
		EndpointPrivateAccess: pulumi.Bool(false),
		EndpointPublicAccess: pulumi.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	ctx.Export("kubeconfig", cluster.Kubeconfig)
	ctx.Export("vpcId", vpc.VpcId)

	return cluster, nil
}