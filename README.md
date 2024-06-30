# Chirpy API
This project was built using a guide from the [boot.dev](https://www.boot.dev) course.
Chirpy is a social network where you can create a user, log in using its credentials and publish posts.

## Goal / Motivation
The goal of the project was to build an API in the Go programming language and learn authentication, authorization and webhooks.

## Overview
Some of the endpoints requires access tokens or refresh tokens to authenticate the user.
Data is stored by reading and writing to a local file that gets created upon application startup.
A made up third-party service called "Polka" is used to experiment how webhooks are made. To verify that it is the Polka that "sends" the webhook, an API-key is used in the auth-header of the request. The API-key is then compared to a secret that is stored locally before proceeding.

