import pulumi
import pulumi_aws as aws
import pulumi_docker_build as docker_build

output_s3_bucket = aws.s3.Bucket(
    "output-s3-bucket",
    bucket="ml-models"
    acl="private"
)

default_vpc = aws.ec2.get_vpc(default=True)
default_subnets = aws.ec2.get_subnets(
    filters=[
        {
            "name": "vpc-id",
            "values": [default_vpc.id],
        }
    ]
)
vpc_security_groups = aws.ec2.get_security_groups(
    filters=[
        {
            "name": "vpc-id",
            "values": [default_vpc.id],
        }
    ]
)

ml_repo = aws.ecr.Repository(
    "ml_ecr",
    name="ml",
    image_scanning_configuration=aws.ecr.RepositoryImageScanningConfigurationArgs(
        scan_on_push=True,
    ),
)

auth_token = aws.ecr.get_authorization_token_output(registry_id=ml_repo.registry_id)
ml_repo_url = ml_repo.repository_url.apply(lambda url: f"{url}:latest")
ml_image = docker_build.Image(
    "ml-image",
    tags=[ml_repo_url],
    context=docker_build.BuildContextArgs(
        location="../catboost",
    ),
    cache_from=[
        docker_build.CacheFromArgs(
            registry=docker_build.CacheFromRegistryArgs(
                ref=ml_repo_url,
            ),
        )
    ],
    cache_to=[
        docker_build.CacheToArgs(
            inline=docker_build.CacheToInlineArgs(),
        )
    ],
    platforms=[docker_build.Platform.LINUX_AMD64],
    push=True,
    registries=[
        docker_build.RegistryArgs(
            address=ml_repo.repository_url,
            password=auth_token.password,
            username=auth_token.user_name,
        )
    ],
)


launch_template = aws.ec2.LaunchTemplate(
    "ml-batch",
    name="ml_batch",
    block_device_mappings=[
        aws.ec2.LaunchTemplateBlockDeviceMappingArgs(
            device_name="/dev/xvda",
            ebs=aws.ec2.LaunchTemplateBlockDeviceMappingEbsArgs(
                volume_size=100,
                volume_type="gp2",
                delete_on_termination="true",
            ),
        )
    ],
    metadata_options=aws.ec2.LaunchTemplateMetadataOptionsArgs(
        http_tokens="required",
    ),
)

service_role = aws.iam.Role(
    "ml-service-role",
    name="ml-service-role",
    assume_role_policy=pulumi.Output.json_dumps(
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"Service": "batch.amazonaws.com"},
                    "Action": "sts:AssumeRole",
                }
            ],
        }
    ),
)

service_role_policy_attachment = aws.iam.RolePolicyAttachment(
    "ml-service-role-attachment",
    role=service_role.name,
    policy_arn="arn:aws:iam::aws:policy/service-role/AWSBatchServiceRole",
)


ecs_instance_role = aws.iam.Role(
    "ecs-instance-role",
    name="ecs-instance-role",
    assume_role_policy=pulumi.Output.json_dumps(
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"Service": "ec2.amazonaws.com"},
                    "Action": "sts:AssumeRole",
                }
            ],
        }
    ),
)

ecs_instance_role_policy_attachment = aws.iam.RolePolicyAttachment(
    "ecs-instance-role-policy-attachment",
    role=ecs_instance_role.name,
    policy_arn="arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role",
)

ecs_instance_profile = aws.iam.InstanceProfile(
    "ecs-instance-profile",
    role=ecs_instance_role.name,
)

log_group = aws.cloudwatch.LogGroup(
    "ml-log-group",
    name="/aws/batch/ml",
    retention_in_days=30,
)

compute_environment = aws.batch.ComputeEnvironment(
    "ml-compute-environment",
    compute_environment_name="ce_ml",
    compute_resources=aws.batch.ComputeEnvironmentComputeResourcesArgs(
        instance_types=["g5.2xlarge"],
        max_vcpus=8,
        min_vcpus=0,
        desired_vcpus=0,
        type="EC2",
        launch_template=aws.batch.ComputeEnvironmentComputeResourcesLaunchTemplateArgs(
            launch_template_id=launch_template.id,
            version="$Latest",
        ),
        security_group_ids=vpc_security_groups.ids,
        subnets=default_subnets.ids,
        instance_role=ecs_instance_profile.arn,
        ec2_configurations=[
            aws.batch.ComputeEnvironmentComputeResourcesEc2ConfigurationArgs(
                image_type="ECS_AL2_NVIDIA",
            )
        ],
    ),
    service_role=service_role.arn,
    type="MANAGED",
    opts=pulumi.ResourceOptions(depends_on=[log_group]),
)

job_queue = aws.batch.JobQueue(
    "ml-job-queue",
    name="ml_job_queue",
    state="ENABLED",
    priority=10,
    compute_environments=[compute_environment.arn],
)

job_role = aws.iam.Role(
    "ml-job-role",
    name="ml-job",
    description="Job role for the ml job",
    assume_role_policy=pulumi.Output.json_dumps(
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Action": "sts:AssumeRole",
                    "Effect": "Allow",
                    "Principal": {"Service": "ecs-tasks.amazonaws.com"},
                }
            ],
        }
    ),
)

job_role_policy = aws.iam.RolePolicy(
    "ml-job-role-policy",
    role=job_role.name,
    policy=output_s3_bucket.arn.apply(
        lambda arn: pulumi.Output.json_dumps(
            {
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Effect": "Allow",
                        "Action": ["s3:PutObject"],
                        "Resource": [f"{arn}/*"],
                    }
                ],
            }
        )
    ),
)

ecs_task_execution_role = aws.iam.Role(
    "ecs_task_execution_role",
    name="ecs-task-execution-role",
    description="Allows ECS and Batch to execute tasks",
    assume_role_policy=pulumi.Output.json_dumps(
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Action": "sts:AssumeRole",
                    "Principal": {"Service": "ecs-tasks.amazonaws.com"},
                }
            ],
        }
    ),
)

ecs_task_execution_attach = aws.iam.RolePolicyAttachment(
    "ecs_task_execution_attach",
    role=ecs_task_execution_role.name,
    policy_arn="arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
)

job_definition = aws.batch.JobDefinition(
    "ml-job-definition",
    name="ml-job-definition",
    type="container",
    platform_capabilities=["EC2"],
    timeout=aws.batch.JobDefinitionTimeoutArgs(
        attempt_duration_seconds=30 * 60,
    ),
    container_properties=pulumi.Output.json_dumps(
        {
            "image": ml_repo.repository_url.apply(lambda url: f"{url}:latest"),
            "stopTimeout": 120,
            "resourceRequirements": [
                {"type": "VCPU", "value": "8"},
                {"type": "MEMORY", "value": "30000"},
                {"type": "GPU", "value": "1"},
            ],
            "linuxParameters": {
                "initProcessEnabled": True,
            },
            "environment": [
                {
                    "name": "AWS_DEFAULTS_MODE",
                    "value": "auto",
                },
                {
                    "name": "OUTPUT_BUCKET",
                    "value": output_s3_bucket.bucket.apply(lambda bucket: bucket),
                },
            ],
            "logConfiguration": {
                "logDriver": "awslogs",
                "options": {
                    "awslogs-group": log_group.name,
                },
            },
            "executionRoleArn": ecs_task_execution_role.arn,
            "jobRoleArn": job_role.arn,
        }
    ),
)

pulumi.export("ml_repo_url", ml_repo.repository_url)
pulumi.export("output_s3_bucket", output_s3_bucket.bucket)
pulumi.export("compute_environment_arn", compute_environment.arn)
pulumi.export("job_queue_arn", job_queue.arn)
pulumi.export("job_definition_arn", job_definition.arn)
