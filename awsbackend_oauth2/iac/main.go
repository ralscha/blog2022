package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/cognito"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		account, err := aws.GetCallerIdentity(ctx)
		if err != nil {
			return err
		}

		region, err := aws.GetRegion(ctx, &aws.GetRegionArgs{})
		if err != nil {
			return err
		}

		clientId, poolId, err := createCognitoUserPool(ctx)
		if err != nil {
			return err
		}

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

		issuer := poolId.ToStringOutput().ApplyT(func(poolId string) string {
			return fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region.Name, poolId)
		}).(pulumi.StringOutput)

		api, err := createApi(ctx, lmbda, clientId, issuer)
		if err != nil {
			return err
		}

		err = createInvokeLambdaPermission(ctx, lmbda, api, region.Name, account.AccountId)
		if err != nil {
			return err
		}

		return nil
	})
}

func createCognitoUserPool(ctx *pulumi.Context) (pulumi.IDOutput, pulumi.IDOutput, error) {
	userPool, err := cognito.NewUserPool(ctx, "pool", &cognito.UserPoolArgs{
		AccountRecoverySetting: &cognito.UserPoolAccountRecoverySettingArgs{
			RecoveryMechanisms: &cognito.UserPoolAccountRecoverySettingRecoveryMechanismArray{
				&cognito.UserPoolAccountRecoverySettingRecoveryMechanismArgs{
					Name:     pulumi.String("verified_email"),
					Priority: pulumi.Int(1),
				},
			},
		},
		AliasAttributes: pulumi.StringArray{
			pulumi.String("preferred_username"),
		},
		AutoVerifiedAttributes: pulumi.StringArray{
			pulumi.String("email"),
		},
		MfaConfiguration: pulumi.String("OFF"),
		Name:             pulumi.String("todo"),
		Schemas: cognito.UserPoolSchemaArray{
			&cognito.UserPoolSchemaArgs{
				AttributeDataType: pulumi.String("String"),
				Mutable:           pulumi.Bool(true),
				Name:              pulumi.String("email"),
				Required:          pulumi.Bool(true),
				StringAttributeConstraints: &cognito.UserPoolSchemaStringAttributeConstraintsArgs{
					MaxLength: pulumi.String("2048"),
					MinLength: pulumi.String("0"),
				},
			},
		},
	})
	if err != nil {
		return pulumi.IDOutput{}, pulumi.IDOutput{}, err
	}

	client, err := cognito.NewUserPoolClient(ctx, "client", &cognito.UserPoolClientArgs{
		CallbackUrls:         pulumi.StringArray{pulumi.String("http://localhost:8100")},
		AccessTokenValidity:  pulumi.Int(1),
		IdTokenValidity:      pulumi.Int(1),
		RefreshTokenValidity: pulumi.Int(30),
		TokenValidityUnits: cognito.UserPoolClientTokenValidityUnitsArgs{
			AccessToken:  pulumi.String("hours"),
			IdToken:      pulumi.String("hours"),
			RefreshToken: pulumi.String("days"),
		},
		EnableTokenRevocation:      pulumi.Bool(true),
		PreventUserExistenceErrors: pulumi.String("ENABLED"),
		AllowedOauthFlows: pulumi.StringArray{
			pulumi.String("code"),
		},
		AllowedOauthFlowsUserPoolClient: pulumi.Bool(true),
		AllowedOauthScopes: pulumi.StringArray{
			pulumi.String("openid"),
		},
		ExplicitAuthFlows: pulumi.StringArray{
			pulumi.String("ALLOW_REFRESH_TOKEN_AUTH"),
		},
		Name: pulumi.String("todo"),
		ReadAttributes: pulumi.StringArray{
			pulumi.String("email"),
		},
		SupportedIdentityProviders: pulumi.StringArray{
			pulumi.String("COGNITO"),
		},
		UserPoolId: userPool.ID(),
		WriteAttributes: pulumi.StringArray{
			pulumi.String("email"),
		},
	})
	if err != nil {
		return pulumi.IDOutput{}, pulumi.IDOutput{}, err
	}

	_, err = cognito.NewUserPoolDomain(ctx, "main", &cognito.UserPoolDomainArgs{
		Domain:     pulumi.String("todo-2021"),
		UserPoolId: userPool.ID(),
	})
	if err != nil {
		return pulumi.IDOutput{}, pulumi.IDOutput{}, err
	}

	ctx.Export("cognito client id", client.ID())
	ctx.Export("cognito user id", userPool.ID())

	return client.ID(), userPool.ID(), nil
}

func createDynamoDbTable(ctx *pulumi.Context) (*dynamodb.Table, error) {
	todoDb, err := dynamodb.NewTable(ctx, "todo-dynamodb-table", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("UserId"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("Id"),
				Type: pulumi.String("S"),
			},
		},
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		HashKey:     pulumi.String("UserId"),
		RangeKey:    pulumi.String("Id"),
	}, pulumi.IgnoreChanges([]string{"read_capacity", "write_capacity"}))
	if err != nil {
		return nil, err
	}
	return todoDb, nil
}

func createInvokeLambdaPermission(ctx *pulumi.Context,
	lmbda *lambda.Function,
	api *apigatewayv2.Api,
	regionName string,
	accountId string) error {
	_, err := lambda.NewPermission(ctx, "todo-apigateway-invoke-lambda-permission", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  lmbda.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", regionName, accountId, api.ID()),
	})
	if err != nil {
		return err
	}

	return nil
}

func createApi(ctx *pulumi.Context, lmbda *lambda.Function, clientId pulumi.IDOutput, issuer pulumi.StringOutput) (*apigatewayv2.Api, error) {
	api, err := apigatewayv2.NewApi(ctx, "todo", &apigatewayv2.ApiArgs{
		ProtocolType: pulumi.String("HTTP"),
		CorsConfiguration: apigatewayv2.ApiCorsConfigurationArgs{
			AllowCredentials: pulumi.Bool(false),
			AllowMethods:     pulumi.StringArray{pulumi.String("GET"), pulumi.String("POST"), pulumi.String("DELETE")},
			AllowOrigins:     pulumi.StringArray{pulumi.String("*")},
			AllowHeaders:     pulumi.StringArray{pulumi.String("content-type"), pulumi.String("authorization")},
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

	authorizer, err := apigatewayv2.NewAuthorizer(ctx, "todo-jwt-authorizer", &apigatewayv2.AuthorizerArgs{
		ApiId:          api.ID(),
		AuthorizerType: pulumi.String("JWT"),
		IdentitySources: pulumi.StringArray{
			pulumi.String("$request.header.Authorization"),
		},
		JwtConfiguration: &apigatewayv2.AuthorizerJwtConfigurationArgs{
			Audiences: pulumi.StringArray{clientId},
			Issuer:    issuer,
		},
	})
	if err != nil {
		return nil, err
	}

	target := integration.ID().ToStringOutput().ApplyT(func(integrationId string) string {
		return "integrations/" + integrationId
	}).(pulumi.StringOutput)

	_, err = apigatewayv2.NewRoute(ctx, "todos-get", &apigatewayv2.RouteArgs{
		ApiId:             api.ID(),
		AuthorizationType: pulumi.String("JWT"),
		AuthorizerId:      authorizer.ID(),
		RouteKey:          pulumi.String("GET /todos"),
		Target:            target,
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewRoute(ctx, "todos-post", &apigatewayv2.RouteArgs{
		ApiId:             api.ID(),
		AuthorizationType: pulumi.String("JWT"),
		AuthorizerId:      authorizer.ID(),
		RouteKey:          pulumi.String("POST /todos"),
		Target:            target,
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewRoute(ctx, "todos-delete", &apigatewayv2.RouteArgs{
		ApiId:             api.ID(),
		AuthorizationType: pulumi.String("JWT"),
		AuthorizerId:      authorizer.ID(),
		RouteKey:          pulumi.String("DELETE /todos/{id}"),
		Target:            target,
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
                      "dynamodb:Query"
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
