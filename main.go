package main

import (
	"context"
	"fmt"
	"maos-cloud-project-api/aws"
	"maos-cloud-project-api/project"
	"maos-cloud-project-api/stack"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)


func main() {
	ensurePlugins()
	// aws.CreateAwsSession()
	routes()

}

func routes() {
	router := mux.NewRouter()

	//project operations
	router.HandleFunc("/project", project.CreateProject).Methods("POST")

	//Stack operations
    router.HandleFunc("/{project_name}/stack", stack.CreateStack).Methods("POST")
	router.HandleFunc("/{project_name}/stack/{stack_name}", stack.DeleteStack).Methods("DELETE")
	router.HandleFunc("/{project_name}/stack/{stack_name}", stack.GetStack).Methods("GET")

	// setup our AWS RESTful routes 
	router.HandleFunc("/{project_name}/stack/{stack_name}/aws/vpc", aws.CreateAwsVPC).Methods("POST")
	router.HandleFunc("/{project_name}/stack/{stack_name}/aws/vpc/{vpc_id}", aws.DeleteVPC).Methods("DELETE")
	router.HandleFunc("/{project_name}/stack/{stack_name}/aws/vpc/{vpc_id}", aws.UpdateVpc).Methods("PUT")

	//define and start the http server
	server := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	// logger      := logrus.New()
	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Errorf("An error occured ", err)

	}
}

/* 
ensure plugins run once before the server boots up
making sure the proper pulumi plugins are installed

*/
func ensurePlugins() {

	fmt.Println("Ensuring all deps are installed")
	ctx := context.Background()
	w, err := auto.NewLocalWorkspace(ctx)

	if err != nil {
		fmt.Printf("Failed to setup and run https server %v\n", err )
		os.Exit(1)
	}

	err = w.InstallPlugin(ctx, "aws", "v3.26.0")
	if err != nil {
		fmt.Printf("Failed to install program plugins: %v\n", err )
		os.Exit(1)
	}
}
