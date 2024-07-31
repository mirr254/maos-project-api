// Description: This file includes the API endpoints that are used to create AWS resources.
package controllers

import (
	"maos-cloud-project-api/helpers/aws"
	"net/http"

	"github.com/gin-gonic/gin"
	awsx "github.com/pulumi/pulumi-awsx/sdk/v2/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/sirupsen/logrus"
)

type VPC struct {
	Name string `json:"project_name" binding:"required"`
	Stack string `json:"stack_name" binding:"required"`
	BastionHostAccessWhitelistCidr string `json:"bastion_host_access_whitelist_cidr" binding:"required"`

}

func CreateVPCEndpoint(c *gin.Context) {

	var vpc VPC
	if err := c.ShouldBindJSON(&vpc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if vpc.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_name must be provided"})	
		return
	}

	var vpcOutput *awsx.Vpc
	var vpcErr error
	var programCreateVpc = func(ctx *pulumi.Context) error {
		vpcOutput, vpcErr = aws_helpers.CreateVPCResource(ctx, vpc.Name)
		return vpcErr
	}

	err := aws_helpers.VPCOperations(programCreateVpc,vpc.Name, vpc.Stack)
	if err != nil {
		logrus.Error("Could not create VPC", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	vpcIDChan := make(chan string)
	vpcOutput.VpcId.ApplyT(func(vpcID string) error {
		vpcIDChan <- vpcID
		return nil
	}).(pulumi.Output).ApplyT(func(_ interface{}) error {
		close(vpcIDChan)
		return nil
	})

	vpcID := <-vpcIDChan

	logrus.Info("VPC created successfully: ", vpcID)
	c.JSON(http.StatusCreated, gin.H{ "vpc_name": vpc.Name, "vpc_id": vpcID})
	return

}
