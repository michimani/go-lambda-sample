go-lambda-sample
---

This is a sample that implements the AWS Lambda function in Go language and deploys it as a container image. The sample function gets a list of S3 Buckets.

## Usage

1. Clone repository

    ```bash
    $ git clone https://github.com/michimani/go-lambda-sample.git
    $ cd go-lambda-sample
    ```

2. Build the image

    ```bash
    $ docker build -t go-lambda-sample .
    ```

## Run at local

1. Download RIE and install it on your local machine

    ```bash
    $ mkdir -p ~/.aws-lambda-rie \
    && curl -Lo ~/.aws-lambda-rie/aws-lambda-rie \
    https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie \
    && chmod +x ~/.aws-lambda-rie/aws-lambda-rie
    ```
2. Run the image

    ```bash
    $ docker run \
    --rm \
    --entrypoint /aws-lambda/aws-lambda-rie \
    -v ~/.aws-lambda-rie:/aws-lambda \
    -p 9000:8080 \
    -e AWS_DEFAULT_REGION="ap-northeast-1" \
    -e AWS_ACCESS_KEY_ID="<your-aws-access-key-id>" \
    -e AWS_SECRET_ACCESS_KEY="<your-aws-secret-acess-key>" \
    go-lambda-sample:latest /main
    ```

3. Invoke function

    ```bash
    $ curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{}'
    ```

## Deploy to Lmabda

1. Build the image

    ```bash
    $ docker build -t go-lambda-sample .
    ```

2. Create ECR repository

    ```bash
    $ aws ecr create-repository \
    --repository-name go-lambda-sample \
    --region ap-northeast-1
    ```

3. Login to ECR

    ```bash
    $ aws ecr get-login-password --region ap-northeast-1 \
    | docker login \
    --username AWS \
    --password-stdin ************.dkr.ecr.ap-northeast-1.amazonaws.com
    ```
    
4. Add image tag

    ```bash
    $ docker tag go-lambda-sample:latest ************.dkr.ecr.ap-northeast-1.amazonaws.com/go-lambda-sample:latest
    ```
    
5. Push to ECR repository

    ```bash
    $ docker push ************.dkr.ecr.ap-northeast-1.amazonaws.com/go-lambda-sample:latest
    ```

6. Create Lambda Function from container image

    ```bash
    $ aws lambda create-function \
    --function-name "go-lambda-sample" \
    --package-type "Image" \
    --code "ImageUri=<your-ecr-repository-uri>" \
    --timeout 30 \
    --role "<iam-role-arn>" \
    --region "<region-code>"
    ```

7. Invoke function

    ```bash
    $ aws lambda invoke \
    --function-name go-lambda-sample \
    --invocation-type RequestResponse \
    --region ap-northeast-1 \
    out \
    && cat out | jq .
    ```
