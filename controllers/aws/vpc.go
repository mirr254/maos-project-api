package aws

import (
	awsec2 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	awsx "github.com/pulumi/pulumi-awsx/sdk/v2/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/sirupsen/logrus"
)

type sgGroup struct {
	group *awsec2.SecurityGroup
}

type jumpbox struct {
	instance *awsec2.Instance
    sgGroup *awsec2.SecurityGroup
}


// CreateVPC creates a new VPC with a public subnet and a private subnet.
func CreateVPCResource(ctx *pulumi.Context, project_name string) ( *awsx.Vpc, error) {

	// Create a new VPC with a public and private subnet.
	vpc, err := awsx.NewVpc(ctx, project_name, &awsx.VpcArgs{
		CidrBlock: pulumi.StringRef("172.16.8.0/24"),
		Tags: pulumi.StringMap{
			"Name": pulumi.String(project_name),
		},
		NumberOfAvailabilityZones: pulumi.IntRef(4),
		SubnetSpecs: []awsx.SubnetSpecArgs{
			{
				Type:     awsx.SubnetTypePublic,
				CidrMask: pulumi.IntRef(22),
				Name:    pulumi.StringRef(project_name+" Public subnet A"),
			},
			{
				Type:     awsx.SubnetTypePublic,
				CidrMask: pulumi.IntRef(22),
				Name:    pulumi.StringRef(project_name+" Public subnet B"),
			},
			{
				Type:     awsx.SubnetTypePrivate,
				CidrMask: pulumi.IntRef(20),
				Name:     pulumi.StringRef(project_name+" Private subnet A"),
			},
			{
				Type:     awsx.SubnetTypePrivate,
				CidrMask: pulumi.IntRef(20),
				Name:     pulumi.StringRef(project_name+" Private subnet N"),
			},
			
		},
		NatGateways: &awsx.NatGatewayConfigurationArgs{
			Strategy: awsx.NatGatewayStrategySingle,
		},
	})
	if err != nil {
		return nil, err
	}

	//create security group
	_, err = createSecurityGroup(ctx, project_name, vpc)
	if err != nil {
		logrus.Error("Could not create security group", err)
		return nil, err
	}

	//create jumpbox
	err = CreateJumpBoxResource(ctx, project_name, vpc)
	if err != nil {
		logrus.Error("Could not create jumpbox", err)
		return nil, err
	}

	
	return vpc, nil

}

//create security group
func createSecurityGroup(ctx *pulumi.Context, project_name string, vpc *awsx.Vpc) ( *sgGroup, error) {
	// Create a new security group.
	group, err := awsec2.NewSecurityGroup(ctx, project_name, &awsec2.SecurityGroupArgs{
		VpcId: vpc.VpcId,
		Ingress: awsec2.SecurityGroupIngressArray{
			&awsec2.SecurityGroupIngressArgs{
				Protocol: pulumi.String("tcp"),
				FromPort: pulumi.Int(80),
				ToPort:   pulumi.Int(80),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0/0"),
				},
			},
			&awsec2.SecurityGroupIngressArgs{
				Protocol: pulumi.String("tcp"),
				FromPort: pulumi.Int(22),
				ToPort:   pulumi.Int(22),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("10.0.0.0/0"), //TDO: Change this to a more secure IP
				},
			},
			&awsec2.SecurityGroupIngressArgs{
				Protocol: pulumi.String("tcp"),
				FromPort: pulumi.Int(443),
				ToPort:   pulumi.Int(443),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),	
				},
			},
		},
		Egress: awsec2.SecurityGroupEgressArray{
			&awsec2.SecurityGroupEgressArgs{
				Protocol: pulumi.String("-1"),
				FromPort: pulumi.Int(0),
				ToPort:   pulumi.Int(0),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
			},
		},
		Tags: pulumi.StringMap{
			"Name": pulumi.String(project_name),
		},
	})

	if err != nil {
		return nil, err
	}
	// ctx.Export("vpcId", pulumi.String(vpcID))
	return &sgGroup{
		group: group,
	}, nil
	
}

//create ec2 instance to be used as a jumpbox
func CreateJumpBoxResource(ctx *pulumi.Context, project_name string, vpc *awsx.Vpc) error {
	// Create a new security group.
	ami, err := awsec2.LookupAmi(ctx, &awsec2.LookupAmiArgs{
		Filters: []awsec2.GetAmiFilter{
			{
				Name:   "name",
				Values: []string{"amzn2-ami-hvm-*"},
			},
		},
			Owners: []string{"amazon"},
			MostRecent: pulumi.BoolRef(true),
		})
		if err != nil {
			logrus.Error("Could not find AMI", err)
			return err
		}

		jumpboxSG, err := awsec2.NewSecurityGroup(ctx, project_name, &awsec2.SecurityGroupArgs{
			VpcId: vpc.VpcId,
			Ingress: awsec2.SecurityGroupIngressArray{
				&awsec2.SecurityGroupIngressArgs{
					Protocol: pulumi.String("tcp"),
					FromPort: pulumi.Int(22),
					ToPort:   pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"), //TDO: Change this to a more secure IP
					},
				},
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String(project_name),
			},
		})
		if err != nil {
			logrus.Error("Could not create security group", err)
			return err
		}

		awsec2.NewInstance(ctx, project_name, &awsec2.InstanceArgs{
			Ami:           pulumi.String(ami.Id),
			InstanceType:  pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{jumpboxSG.ID()},
			SubnetId: vpc.PublicSubnetIds.Index(pulumi.Int(0)),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(project_name),
			},
		})
		if err != nil {
			logrus.Error("Could not create instance", err)
			return err
		}
		return nil

}

