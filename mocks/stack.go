package mocks

// import (
// 	"context"

// 	"github.com/pulumi/pulumi/sdk/v3/go/auto"
// )

// type StackCreator interface {
// 	NewStackInlineSource(ctx context.Context, stackName, projectName string, program auto.ProgramFunc, opts ...auto.LocalWorkspaceOption) (auto.Stack, error)
// }

// type PulumiStackCreator struct{}

// func (psc *PulumiStackCreator) NewStackInlineSource(ctx context.Context, stackName, projectName string, program auto.Program, opts ...auto.LocalWorkspaceOption) (auto.Stack, error) {
// 	return auto.NewStackInlineSource(ctx, stackName, projectName, program, opts...)
// }
