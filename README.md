# Ascendio 🏊
[![GitHub license](https://img.shields.io/github/license/vickyphang/ascendio)](https://github.com/vickyphang/ascendio/blob/main/LICENSE)
![GitHub stars](https://img.shields.io/github/stars/vickyphang/ascendio)

### Github Apps Integration for Hyperback Project
<p align="center"> <img src="images/logo.png"> </p>

## Overview
This project generate a login button and install the github app to a repository also generate access token + jwt token for API interaction

A GitHub App is a type of integration that you can build to interact with and extend the functionality of GitHub. You can build a GitHub App to provide flexibility and reduce friction in your processes, without needing to sign in a user or create a service account.

<p align="center"> <img src="images/flowchart.png"> </p>

## Setup a Github Apps
#### 1. Register github apps
- Go to GitHub's `Developer Settings` and create a new `GitHub App`: https://github.com/settings/apps.

#### 2. Fill in the necessary details:
- `GitHub App Name`: Choose a unique name for your app.
- `Homepage URL`: Set it to your server's URL (for development, http://localhost:8080).
- `Callback URL`: Set it to http://localhost:8080/callback.
- `Webhook URL`: (Optional: you can turn it off) Set it to http://localhost:8080/webhook. Add **webhook secret** to improve security
- `Permissions`: Set the permissions your app will need.
- `Events`: Subscribe to Installation and InstallationRepositories events.
- `Where can this GitHub App be installed?`: Choose `Only on this account` for development purpose

#### 3. Save the credentials `App ID` and `Client ID`
<p align="center"> <img src="images/github-apps.png"> </p>

#### 4. Generate a new client secret and store it
<p align="center"> <img src="images/client-secret.png"> </p>

#### 5. Generate a private key and save it as private-key.pem
<p align="center"> <img src="images/private-key.png"> </p>


## Setup Golang Application
Clone this repository into your server
```bash
git clone https://github.com/vickyphang/ascendio.git

cd ascendio
```

Make sure you have `docker` installed. Build the application from `Dockerfile`
```bash
docker build -t golang-app:v1 .
```

Copy `private-key.pem` to this folder. Run the application
```bash
docker run -p 8080:8080 --name golang-app /
    -e CLIENT_ID="<client-id>" /
    -e CLIENT_SECRET="<client-secret>" /
    -e APP_ID="<app-id>" /
    -e WEBHOOK_SECRET="<webhook-secret>" /
    -e PRIVATE_KEY_PEM="<base64-encoded-key-pem>" /
    golang-app:v1
```

Open browser and access: `http://host-ip:8080/login`
<p align="center"> <img src="images/login.png"> </p>

When authorized, Github redirects back to your callback url with `Authorization Code`. Exchange that **Authorization Code** and you will get `Access Token`(and jwt, but you have to uncomment the code). **Access Token** expires in **1 hour**
<p align="center"> <img src="images/callback.png"> </p>

In order to use the `access token` or `jwt token` you must install the Github App to your repository with appropriate permission. To install github apps, access: `http://host-ip:8080/install`
<p align="center"> <img src="images/install.png"> </p>

## Verify
After installing Github App to your selected private repository, try to clone these private repo with access token you got
```bash
git clone https://x-access-token:TOKEN@github.com/owner/private-repo.git
```

You also can use the access token to make an API request as an installation
```
curl --location 'https://api.github.com/repos/vickyphang/ascendio/commits' \
--header 'Authorization: Bearer <access-token>'
```

> NOTE: Some Github REST API endpoints do not accept installation access tokens