---
title: "NGINX Plus Image and JWT Requirement"
weight: 300
toc: true
docs: "DOCS-000"
---

## Overview

NGINX Gateway Fabric with NGINX Plus requires a valid JSON Web Token (JWT) to download the container image from the F5 registry. In addition, starting with version 1.5.0, this JWT token is also required to run NGINX Plus.

This requirement is part of F5’s broader licensing program and aligns with industry best practices. The JWT will streamline subscription renewals and usage reporting, helping you manage your NGINX Plus subscription more efficiently. The [telemetry](#telemetry) data we collect helps us improve our products and services to better meet your needs.

The JWT is required for validating your subscription and reporting telemetry data. For environments connected to the internet, telemetry is automatically sent to F5’s licensing endpoint. In offline environments, telemetry is routed through [NGINX Instance Manager](https://docs.nginx.com/nginx-management-suite/nim/). Usage is reported every hour and on startup whenever NGINX is reloaded.

## Setting up the JWT

The JWT needs to be configured before deploying NGINX Gateway Fabric. We'll store the JWT in two Kubernetes Secrets. One will be used for downloading the NGINX Plus container image, and the other for running NGINX Plus.

{{< include "installation/jwt-password-note.md" >}}

### Download the JWT from MyF5

{{<include "installation/nginx-plus/download-jwt.md" >}}

### Docker Registry Secret

{{<include "installation/nginx-plus/docker-registry-secret.md" >}}

Provide the name of this Secret when installing NGINX Gateway Fabric:

{{<tabs name="docker-secret-install">}}

{{%tab name="Helm"%}}

Specify the Secret name using the `nginxGateway.serviceAccount.imagePullSecret` or `nginxGateway.serviceAccount.imagePullSecrets` helm value.

{{% /tab %}}

{{%tab name="Manifests"%}}

Specify the Secret name in the `imagePullSecrets` field of the `nginx-gateway` ServiceAccount.

{{% /tab %}}

{{</tabs>}}

### NGINX Plus Secret

{{<include "installation/nginx-plus/nginx-plus-secret.md" >}}

Provide the name of this Secret when installing NGINX Gateway Fabric:

{{<tabs name="plus-secret-install">}}

{{%tab name="Helm"%}}

Specify the Secret name using the `nginx.usage.secretName` helm value.

{{% /tab %}}

{{%tab name="Manifests"%}}

Specify the Secret name in the `--usage-report-secret` command-line flag on the `nginx-gateway` container.

{{% /tab %}}

{{</tabs>}}

[Additional configuration options](#flags) can be set to further configure usage reporting.

<br>

**Once these two Secrets are created, you can now [install NGINX Gateway Fabric]({{< relref "installation/installing-ngf" >}}).**

## Installation flags for configuring usage reporting {#flags}

When installing NGINX Gateway Fabric, the following flags can be specified to configure usage reporting to fit your needs:

If using Helm, the `nginx.usage` values should be set as necessary:

- `secretName` should be the `name` of the JWT Secret you created. Using our example, it would be `nplus-license`. This field is required.
- `endpoint` is the endpoint to send the telemetry data to. This is optional, and by default is `product.connect.nginx.com`.
- `resolver` is the resolver domain name or IP address with optional port for resolving the endpoint. This is optional.

If using manifests, the following command-line options should be set as necessary on the `nginx-gateway` container:

- `--usage-report-secret` should be the `name` of the JWT Secret you created. Using our example, it would be `nplus-license`. This field is required.
- `--usage-report-endpoint` is the endpoint to send the telemetry data to. This is optional, and by default is `product.connect.nginx.com`.
- `--usage-report-resolver` is the resolver domain name or IP address with optional port for resolving the endpoint. This is optional.

## What’s reported and how it’s protected {#telemetry}

NGINX Plus reports the following data every hour by default:

- **NGINX version and status**: The version of NGINX Plus running on the instance.
- **Instance UUID**: A unique identifier for each NGINX Plus instance.
- **Traffic data**:
  - **Bytes received from and sent to clients**: HTTP and stream traffic volume between clients and NGINX Plus.
  - **Bytes received from and sent to upstreams**: HTTP and stream traffic volume between NGINX Plus and upstream servers.
  - **Client connections**: The number of accepted client connections (HTTP and stream traffic).
  - **Requests handled**: The total number of HTTP requests processed.
- **NGINX uptime**: The number of reloads and worker connections during uptime.
- **Usage report timestamps**: Start and end times for each usage report.
- **Kubernetes node details**: Information about Kubernetes nodes.

### Security and privacy of reported data

All communication between your NGINX Plus instances, NGINX Instance Manager, and F5’s licensing endpoint (`product.connect.nginx.com`) is protected using **SSL/TLS** encryption.

Only **operational metrics** are reported — no **personally identifiable information (PII)** or **sensitive customer data** is transmitted.

## Pulling an image for local use

To pull an image for local use, use this command:

```shell
docker login private-registry.nginx.com --username=<JWT Token> --password=none
```

Replace the contents of `<JWT Token>` with the contents of the JWT token itself.

You can then pull the image:

```shell
docker pull private-registry.nginx.com/nginx-gateway-fabric/nginx-plus:1.4.0
```

Once you have successfully pulled the image, you can tag it as needed, then push it to a different container registry.

## Alternative installation options

There are alternative ways to get an NGINX Plus image for NGINX Gateway Fabric:

- [Build the Gateway Fabric image]({{<relref "installation/building-the-images.md">}}) describes how to use the source code with an NGINX Plus subscription certificate and key to build an image.
