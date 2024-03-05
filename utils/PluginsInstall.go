package utils

import (
	"os"
	"context"

	auto "github.com/pulumi/pulumi/sdk/v3/go/auto"
	log  "github.com/sirupsen/logrus"

)

/*
ensure plugins run once before the server boots up
making sure the proper pulumi plugins are installed
*/
func EnsurePlugins() {

	log.Info("Ensuring all deps are installed")
	ctx := context.Background()
	w, err := auto.NewLocalWorkspace(ctx)

	if err != nil {
		log.Printf("Failed to setup and run https server %v\n", err)
		os.Exit(1)
	}

	err = w.InstallPlugin(ctx, "aws", "v3.26.0")
	if err != nil {
		log.Printf("Failed to install aws plugin: %v\n", err)
		os.Exit(1)
	}
}