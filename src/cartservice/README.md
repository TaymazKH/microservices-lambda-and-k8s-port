## Language

This service was translated to Golang to simply the porting process and avoid complications of using C#.

## Setting Up a Redis Cache

This service utilizes a cache system to store carts. Two cache systems are implemented for this service: a Redis cache
and an in-memory object map. Due to the serverless nature of Lambda, the in-memory map can't be used in this deployment
method. Therefore, a Redis cache needs to be set up.

For the purpose of this tutorial, we will use AWS ElastiCache with a Redis engine. Create a custom cluster. For
simplicity, disable settings such as backups and replicas. Set up one node and set its type to `cache.t3.micro`, as this
was the most cost-efficient option at the time of this document's writing. You may want to take a look at the pricing
policy and choose another node type. Create a VPC (or use an existing one) and specify the subnets. Disable encryption
for simplicity. Now, wait until your cache creation is finalized.

If you haven't explicitly set the permissions of your Lambda functions, they will use the most basic execution role.
This role has only the `AWSLambdaBasicExecutionRole` policy which does not have the required permissions for VPC access.
You need to use a role that grants these permissions. Go to AWS IAM, then the roles sections, and create a new role.
Select "AWS service" as the trusted entity type and "Lambda" as the use case. Next, select
the `AWSLambdaVPCAccessExecutionRole` policy. Finalize your role creation with a name.

Go to your Lambda function, and go to the permissions sections under the configurations tab. Edit the execution role and
choose the new role. Next, go to the VPC section under the same tab. Edit it and choose the same VPC, subnets, and
security groups as the Redis cache. Your Lambda function is now ready to use your ElastiCache Redis cluster.

## Environment Variables

This is the list of unique environment variables this service uses.

Required:

- `REDIS_ADDR`: only in Lambda deployment.

Optional:

- `REDIS_PASS`: only if the Redis cache uses encryption in transit.
