package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		todoDb, err := createDynamoDbTable(ctx)
		if err != nil {
			return err
		}

		role, err := createIamRoleForLambda(ctx, todoDb)
		if err != nil {
			return err
		}

		lmbda, err := createLambda(ctx, role, todoDb)
		if err != nil {
			return err
		}

		api, err := createApi(ctx, lmbda)
		if err != nil {
			return err
		}

		err = createInvokeLambdaPermission(ctx, lmbda, api)
		if err != nil {
			return err
		}

		return nil
	})
}

func createDynamoDbTable(ctx *pulumi.Context) (*dynamodb.Table, error) {
	todoDb, err := dynamodb.NewTable(ctx, "todo-dynamodb-table", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("Id"),
				Type: pulumi.String("S"),
			},
		},
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		HashKey:     pulumi.String("Id"),
	}, pulumi.IgnoreChanges([]string{"read_capacity", "write_capacity"}))
	if err != nil {
		return nil, err
	}
	return todoDb, nil
}

func createInvokeLambdaPermission(ctx *pulumi.Context, lmbda *lambda.Function, api *apigatewayv2.Api) error {
	account, err := aws.GetCallerIdentity(ctx)
	if err != nil {
		return err
	}

	region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
	if err != nil {
		return err
	}

	_, err = lambda.NewPermission(ctx, "todo-apigateway-invoke-lambda-permission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  lmbda.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", region.Name, account.AccountId, api.ID()),
	})
	if err != nil {
		return err
	}

	return nil
}

func createApi(ctx *pulumi.Context, lmbda *lambda.Function) (*apigatewayv2.Api, error) {
	api, err := apigatewayv2.NewApi(ctx, "todo", &apigatewayv2.ApiArgs{
		ProtocolType: pulumi.String("HTTP"),
		CorsConfiguration: apigatewayv2.ApiCorsConfigurationArgs{
			AllowCredentials: pulumi.Bool(false),
			AllowMethods:     pulumi.StringArray{pulumi.String("GET"), pulumi.String("POST"), pulumi.String("DELETE")},
			AllowOrigins:     pulumi.StringArray{pulumi.String("*")},
			AllowHeaders:     pulumi.StringArray{pulumi.String("content-type")},
			MaxAge:           pulumi.Int(7200),
		},
	})
	if err != nil {
		return nil, err
	}

	integration, err := apigatewayv2.NewIntegration(ctx, "todo", &apigatewayv2.IntegrationArgs{
		ApiId:                api.ID(),
		IntegrationType:      pulumi.String("AWS_PROXY"),
		IntegrationUri:       lmbda.Arn,
		IntegrationMethod:    pulumi.String("POST"),
		PayloadFormatVersion: pulumi.String("2.0"),
		TimeoutMilliseconds:  pulumi.Int(3000),
	})

	if err != nil {
		return nil, err
	}

	target := integration.ID().ToStringOutput().ApplyT(func(integrationId string) string {
		return "integrations/" + integrationId
	}).(pulumi.StringOutput)

	_, err = apigatewayv2.NewRoute(ctx, "todos-get", &apigatewayv2.RouteArgs{
		ApiId:    api.ID(),
		RouteKey: pulumi.String("GET /todos"),
		Target:   target,
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewRoute(ctx, "todos-post", &apigatewayv2.RouteArgs{
		ApiId:    api.ID(),
		RouteKey: pulumi.String("POST /todos"),
		Target:   target,
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewRoute(ctx, "todos-delete", &apigatewayv2.RouteArgs{
		ApiId:    api.ID(),
		RouteKey: pulumi.String("DELETE /todos/{id}"),
		Target:   target,
	})
	if err != nil {
		return nil, err
	}

	stage, err := apigatewayv2.NewStage(ctx, "todo-v1-stage", &apigatewayv2.StageArgs{
		ApiId:       api.ID(),
		AutoDeploy:  pulumi.Bool(true),
		Description: pulumi.String("Todo API V1 Stage"),
		Name:        pulumi.String("v1"),
		DefaultRouteSettings: &apigatewayv2.StageDefaultRouteSettingsArgs{
			ThrottlingBurstLimit: pulumi.Int(10),       // maximum number of concurrent requests
			ThrottlingRateLimit:  pulumi.Float64(50.0), // max requests per seconds
		},
	})
	if err != nil {
		return nil, err
	}

	ctx.Export("production stage endpoint", stage.InvokeUrl)

	return api, nil
}

func createIamRoleForLambda(ctx *pulumi.Context, table *dynamodb.Table) (*iam.Role, error) {

	dynamoDbPermissions := table.Arn.ApplyT(func(arn string) string {
		return fmt.Sprintf(`{
          "Version": "2012-10-17",
          "Statement": [
              {
                  "Effect": "Allow",
                  "Action": [
                      "dynamodb:PutItem",
                      "dynamodb:DeleteItem",
                      "dynamodb:Scan"
                  ],
                  "Resource": "%s"
              }
          ]
      }`, arn)
	}).(pulumi.StringOutput)

	role, err := iam.NewRole(ctx, "todo-lambda-exec-role", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
							"Version": "2012-10-17",
							"Statement": [{
								"Sid": "",
								"Effect": "Allow",
								"Principal": {
									"Service": "lambda.amazonaws.com"
								},
								"Action": "sts:AssumeRole"
							}]
						}`),
		InlinePolicies: iam.RoleInlinePolicyArray{iam.RoleInlinePolicyArgs{
			Name:   pulumi.String("dynamodb"),
			Policy: dynamoDbPermissions,
		}},
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, "todo-lambda-exec", &iam.RolePolicyAttachmentArgs{
		Role:      role.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func createLambda(ctx *pulumi.Context, role *iam.Role, todoDb *dynamodb.Table) (*lambda.Function, error) {
	logGroup, err := cloudwatch.NewLogGroup(ctx, "todo", &cloudwatch.LogGroupArgs{
		Name:            pulumi.String("/aws/lambda/todo"),
		RetentionInDays: pulumi.Int(30),
	})
	if err != nil {
		return nil, err
	}

	codeArchive := pulumi.NewAssetArchive(map[string]interface{}{
		"bootstrap": pulumi.NewFileAsset("../lambda/main"),
	})

	args := &lambda.FunctionArgs{
		Runtime:       pulumi.String("provided.al2"),
		Handler:       pulumi.String("bootstrap"),
		Code:          codeArchive,
		MemorySize:    pulumi.Int(128),
		Name:          pulumi.String("todo"),
		Publish:       pulumi.Bool(false),
		Role:          role.Arn,
		Timeout:       pulumi.Int(3),
		Architectures: pulumi.StringArray{pulumi.String("arm64")},
		Environment:   &lambda.FunctionEnvironmentArgs{Variables: pulumi.StringMap{"TABLE_NAME": todoDb.Name}},
	}

	function, err := lambda.NewFunction(
		ctx,
		"todo",
		args,
		pulumi.DependsOn([]pulumi.Resource{role, logGroup}),
	)
	if err != nil {
		return nil, err
	}
	return function, nil
}
