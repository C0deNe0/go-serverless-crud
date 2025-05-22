# Serverless – Overview & Practical Guide

This document provides a simple explanation of serverless architecture, its benefits, typical use cases, and known limitations.

---

## 🌐 What is Serverless?

**Serverless** doesn't mean there are no servers—it means you don’t have to manage them.

In serverless computing:

* You write code and deploy it.
* The cloud provider (like AWS, GCP, Azure) handles provisioning, scaling, and maintenance.
* You only pay for the compute time used during actual execution.

Serverless is often used with **FaaS (Function as a Service)** like AWS Lambda.

---

## ✅ Benefits of Serverless

* **No server management** – No need to provision or manage infrastructure.
* **Automatic scaling** – Handles thousands of requests without any configuration.
* **Pay-as-you-go** – You only pay when your code runs.
* **Faster time to market** – Easier to prototype and ship features.

---

## ❌ Limitations / Challenges of Serverless

* **Cold starts** – Functions may have startup latency if idle for a while.
* **Timeout limits** – Serverless functions are short-lived (e.g., 15 min max for AWS Lambda).
* **State management** – Stateless functions require external storage (DBs, caches, etc).
* **Complex debugging** – Harder to debug locally compared to traditional apps.
* **Vendor lock-in** – Deeply tied to specific cloud platforms (e.g., AWS SAM only for AWS).

---

## 📌 When to Use Serverless

Serverless is ideal for:

* APIs and microservices
* Backend for web/mobile apps
* Scheduled tasks (cron jobs)
* Event-driven systems (e.g., S3 file uploads, DB changes)
* Lightweight backend services

---

## 🧪 Example Use Case

> Build a CRUD API for employee records using:

* FaaS for compute
* API Gateway for HTTP endpoints
* DynamoDB or similar for data storage

---

## 📚 References

* [AWS Lambda](https://docs.aws.amazon.com/lambda/)
* [Serverless Use Cases](https://aws.amazon.com/serverless/)
* [FaaS Concepts](https://martinfowler.com/articles/serverless.html)

---

# AWS SAM Serverless 

This repository contains a Go-based serverless CRUD API built using the **AWS Serverless Application Model (SAM)**. The project utilizes AWS Lambda and API Gateway, and can be deployed directly to the AWS Cloud.

---

## 🧩 What is AWS SAM?

**AWS SAM (Serverless Application Model)** is an open-source framework developed by AWS to simplify the process of building, testing, and deploying serverless applications.

With SAM, you can:

* Define serverless resources using a YAML-based `template.yaml`.
* Run AWS services locally for development.
* Build and deploy applications quickly using the `sam` CLI.

SAM simplifies infrastructure as code for AWS services like:

* **AWS Lambda**
* **API Gateway**
* **S3**

> ⚠️ **Note**: SAM is AWS-specific. For multi-cloud, use tools like Serverless Framework, Pulumi, or Terraform.

---

## ⚙️ Step 1: Install AWS CLI

### ✅ Linux (Ubuntu/Debian):

```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
```

### ✅ macOS (Homebrew):

```bash
brew install awscli
```

### ✅ Windows:

* Download MSI: [https://awscli.amazonaws.com/AWSCLIV2.msi](https://awscli.amazonaws.com/AWSCLIV2.msi)
* Run installer

**Check installation:**

```bash
aws --version
```

---

## ⚙️ Step 2: Install AWS SAM CLI

### ✅ Linux:

```bash
curl -Lo sam-installation.sh https://github.com/aws/aws-sam-cli/releases/latest/download/install
chmod +x sam-installation.sh
sudo ./sam-installation.sh --update
```

### ✅ macOS:

```bash
brew tap aws/tap
brew install aws-sam-cli
```

### ✅ Windows:

* Download: [https://github.com/aws/aws-sam-cli/releases/latest](https://github.com/aws/aws-sam-cli/releases/latest)
* Run installer

**Check installation:**

```bash
sam --version
```

---

## 🔐 Step 3: Configure AWS CLI Credentials

```bash
aws configure
```

Enter:

* AWS Access Key ID
* AWS Secret Access Key
* Default region: `ap-south-1` (or your preferred region)
* Output format: `json`

📚 [AWS Regions](https://docs.aws.amazon.com/global-infrastructure/latest/regions/aws-regions.html)

---

## 📁 Step 4: Initialize Your SAM Project

```bash
sam init
```

Follow prompts:

```text
Which template source would you like to use?
> 1 - AWS Quick Start Templates

Choose an AWS Quick Start application template
> 1 - Hello World Example

Which runtime would you like to use?
> go1.x

Project name:
> helloworld-go
```

---

## 🪣 Step 5: Create an S3 Bucket for Deployment Artifacts

Before deployment, you need a unique S3 bucket:

```bash
aws s3 mb s3://<your-unique-bucket-name> --region <your-region>
```

**Check bucket:**

```bash
aws s3 ls
```

---

## ⚙️ Step 6: Configure `samconfig.toml`

Example config:

```toml
[default.global.parameters]
stack_name = "your_stack_name"

[default.deploy.parameters]
capabilities = "CAPABILITY_IAM"
confirm_changeset = false
s3_bucket = "your_s3_bucket_name"
region = "your_aws_region"
disable_rollback = true

[default.package.parameters]
resolve_s3 = false
```

> ℹ️ This file allows you to skip `sam deploy --guided` and use short deploy commands.

---

## 📄 `template.yaml`

This is the AWS SAM template defining your serverless app.

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Description: >
  serverless-crud
  Sample SAM Template for serverless-crud

Globals:
  Function:
    Timeout: 5
    MemorySize: 128

Resources:
  NaveenFunction:
    Type: AWS::Serverless::Function
    Metadata:
      BuildMethod: go1.x
    Properties:
      CodeUri: src/
      Handler: bootstrap
      Runtime: provided.al2023
      Architectures:
        - x86_64
      Events:
        CatchAll:
          Type: Api
          Properties:
            Path: /{proxy+}
            Method: ANY
      Environment:
        Variables:
          PARAM1: VALUE
```

---

## 🔨 Step 7: Build Your Application

```bash
sam build
```

---

## 🚀 Step 8: Deploy to AWS

### ✅ Recommended: Guided deploy

```bash
sam deploy --guided
```

Follow prompts:

* Stack name
* Region
* Confirm changeset
* Allow role creation
* S3 bucket
* Save config to `samconfig.toml`

### ✅ Shorthand deploy:

```bash
sam deploy --stack-name <your_stack_name> --region <your-region>
```

---

## 🧪 Optional: Local Development (Your Exploration)

Run API locally:

```bash
sam local start-api
```

Invoke Lambda locally:

```bash
sam local invoke "NaveenFunction"
```

---

## 📚 References

* [AWS SAM GitHub](https://github.com/aws/aws-sam-cli)
* [YouTube Guide (Recommended)](https://www.youtube.com/watch?v=IA90BTozdow)
* [AWS Docs](https://docs.aws.amazon.com/serverless-application-model/)

---

> 🛠 Maintained by: \[Vithsutra-serverless Team]
>
> Happy shipping with Serverless! 🚀
