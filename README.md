# _Gamil_ OMG Microservice

[![Open Microservice Guide](https://img.shields.io/badge/OMG%20Enabled-👍-green.svg?)](https://microservice.guide)
[![Build Status](https://travis-ci.com/omg-services/gmail.svg?branch=master)](https://travis-ci.com/omg-services/gmail)
[![codecov](https://codecov.io/gh/omg-services/gmail/branch/master/graph/badge.svg)](https://codecov.io/gh/omg-services/gmail)

An OMG service for Gmail, this service uses the gmail API and perform gmail operations.

## Direct usage in [Storyscript](https://storyscript.io/):

##### Authorization
```coffee
gmail authorization
```
##### Access Token
```coffee
gmail accessToken authorizationCode:'authorization code generated from URL of authorization'
```
##### Send Mail
```coffee
gmail sendMail accessToken:'access token' from:'abc@example.com' to:'["xyz@example.com",mnop@example.com]' subject:'mail subject' body:'mail body' accessToken:'accessToken'
```
##### Receive Mail
```coffee
gmail subscribe receive mail userId:'email address' accessToken:'access token'
```
Curious to [learn more](https://docs.storyscript.io/)?

✨🍰✨

## Usage with [OMG CLI](https://www.npmjs.com/package/omg)

##### Authorization
```shell
omg run authorization -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Access Token
```shell
$ omg run accessToken -a authorizationCode=<AUTHORIZATION_CODE> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Send Mail
```shell
$ omg run sendMail -a accessToken=<ACCESS_TOKEN> -a from=<SENDER_EMAIL_ADDRESS> -a to=[LIST_OF_RECEIVER_EMAIL_ADDRESS] -a subject=<MAIL_SUBJECT> -a body=<MESSAGE_BODY> -a accessToken=<ACCESS_TOKEN> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Receive Mail
```shell
$ omg run subscribe receive mail -a userId=<EMAIL_ADDRESS> -a accessToken=<ACCESS_TOKEN> -e CREDENTIAL_JSON==<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```

**Note**: the OMG CLI requires [Docker](https://docs.docker.com/install/) to be installed.

## License
[MIT License](https://github.com/heaptracetechnology/gmail/blob/master/LICENSE).
