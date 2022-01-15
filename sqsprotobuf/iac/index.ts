import * as aws from "@pulumi/aws";

const queue = new aws.sqs.Queue("queue", {
  delaySeconds: 0,
  maxMessageSize: 262144,
  messageRetentionSeconds: 180,
  receiveWaitTimeSeconds: 20,
  visibilityTimeoutSeconds: 120
});

const bucket = new aws.s3.Bucket("messages", {
  acl: "private",
  lifecycleRules: [{
    enabled: true,
    expiration: {
      days: 1
    }
  }]
});

export const queueArn = queue.url
export const bucketName = bucket.id
