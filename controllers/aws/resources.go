// Description: This file includes the API endpoints that are used to create AWS resources.
package controllers

import (
	aws_helpers "maos-cloud-project-api/helpers/aws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/sirupsen/logrus"
)

type VPC struct {
	Name string `json:"project_name" binding:"required"`
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
	
	ctx := &pulumi.Context{}

	vpcResource, err := aws_helpers.CreateVPCResource(ctx, vpc.Name)
	if err != nil {
		logrus.Error("Could not create VPC", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"vpc": vpcResource})
	return

}
