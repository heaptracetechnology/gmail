# _Gamil_ OMG Microservice

[![Open Microservice Guide](https://img.shields.io/badge/OMG%20Enabled-👍-green.svg?)](https://microservice.guide)
[![Build Status](https://travis-ci.com/omg-services/gmail.svg?branch=master)](https://travis-ci.com/omg-services/gmail)
[![codecov](https://codecov.io/gh/omg-services/gmail/branch/master/graph/badge.svg)](https://codecov.io/gh/omg-services/gmail)

An OMG service for Gmail, this service uses the gmail API and perform gmail operations.

## Direct usage in [Storyscript](https://storyscript.io/):

## NOTE:
use scope "https://mail.google.com/" for send, receive and labels operation while authorization and for "filter operations" use scope "https://www.googleapis.com/auth/gmail.settings.sharing".

##### Authorization
```coffee
gmail authorization scope:'https://mail.google.com/'
```
##### Access Token
```coffee
gmail accessToken authorizationCode:'authorization code generated from URL of authorization'
```
##### Refresh Token
```coffee
gmail refreshToken token:'{"access_token": "access token","token_type": "Bearer","refresh_token": "refresh token","expiry": "2019-10-04T15:57:07.922121141Z"}'
```
##### Send Mail
```coffee
gmail sendMail accessToken:'access token' userId:'abc@example.com' to:'["xyz@example.com",mnop@example.com]' subject:'mail subject' body:'mail body'
```
##### Create Label
```coffee
omg run createLabel accessToken:'access token' userId:'abc@example.com' name:'OMG Label' backgroundColor:'#fce8b3' textColor:'#d5ae49' labelListVisibility:'labelShow' messageListVisibility:'show' 
```
##### Delete Label:
```coffee
omg run deleteLabel accessToken:'access token' userId:'abc@example.com' labelId:'label Id'
```
##### Patch Label:
```coffee
omg run patchLabel accessToken:'access token' userId:'abc@example.com' name:'OMG Label' backgroundColor:'#fce8b3' textColor:'#d5ae49' labelListVisibility:'labelShow' messageListVisibility:'show' 
```
##### List Label
```coffee
omg run labelList accessToken:'access token' userId:'abc@example.com' 
```
##### Create Filter
```coffee
omg run createFilter accessToken:'access token' userId:'abc@example.com'  addLabelId='Label Id' removeLabelId='Label Id' excludeChats=false from:'abc@example.com' to:'xyz@example.com' subject:'Mail subject' hasAttachment=false negatedQuery:'Negated query for filter' query:'query for filter' size=1 sizeComparison:smaller 
```
##### Delete Filter
```coffee
omg run deleteFilter accessToken:'access token' userId:'abc@example.com' filterId:'filter Id'
```
##### List Filter
```coffee
omg run filterList accessToken:'access token' userId:'abc@example.com' 
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
omg run authorization -a scope=<OPERATION_SCOPE> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Access Token
```shell
$ omg run accessToken -a authorizationCode=<AUTHORIZATION_CODE> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Refresh Token
```shell
$ omg run refreshToken -a token=<TOKEN_MAPPED_OBJECT> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Send Mail
```shell
$ omg run sendMail -a accessToken=<ACCESS_TOKEN> -a userId=<SENDER_EMAIL_ADDRESS> -a to=[LIST_OF_RECEIVER_EMAIL_ADDRESS] -a subject=<MAIL_SUBJECT> -a body=<MESSAGE_BODY> -a accessToken=<ACCESS_TOKEN> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Create Label
```shell
omg run createLabel -a accessToken=<ACCESS_TOKEN> -a userId=<EMAIL_ADDRESS> -a name=<LABEL_NAME> -a backgroundColor=<Label_BACKGROUND_COLOR> -a textColor=<Label_TEXT_COLOR> -a labelListVisibility=<LABEL_LIST_VISIBILITY> -a messageListVisibility=<M_LIST_VISIBILITY> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Delete Label:
```shell
omg run deleteLabel -a accessToken=<ACCESS_TOKEN> -a userId=<EMAIL_ADDRESS> -a labelId=<LABEL_ID> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Patch Label:
```shell
omg run patchLabel -a accessToken=<ACCESS_TOKEN> -a userId=<EMAIL_ADDRESS> -a name=<LABEL_NAME> -a backgroundColor=<Label_BACKGROUND_COLOR> -a textColor=<Label_TEXT_COLOR> -a labelListVisibility=<LABEL_LIST_VISIBILITY> -a messageListVisibility=<M_LIST_VISIBILITY> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### List Label
```shell
omg run labelList -a accessToken=<ACCESS_TOKEN> -a userId=<EMAIL_ADDRESS> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Create Filter
```shell
omg run createFilter -a accessToken=<ACCESS_TOKEN> -a userId=<EMAIL_ADDRESS> -a addLabelId=<ADD_LABEL_ID> -a removeLabelId=<REMOVE_LABEL_ID> -a excludeChats=<EXCLUDE_CHATS> -a from=<SENDER_EMAIl_ADDRESS> -a to=<RECEIVER_EMAIL_ADDRESS> -a subject=<MAIL_SUBJECT> -a hasAttachment=<HAS_ATTACHMENT> -a negatedQuery=<NEGATED_QUERY> -a query=<QUERY> -a size=<SIZE_OF_MESSAGE> -a sizeComparison=<SIZE_COMPARISON> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Delete Filter
```shell
omg run deleteFilter -a accessToken=<ACCESS_TOKEN> -a userId=<EMAIL_ADDRESS> -a filterId=<FILTER_ID> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### List Filter
```shell
omg run filterList -a accessToken=<ACCESS_TOKEN> -a userId=<EMAIL_ADDRESS> -e CREDENTIAL_JSON=<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```
##### Receive Mail
```shell
$ omg run subscribe receive mail -a userId=<EMAIL_ADDRESS> -a accessToken=<ACCESS_TOKEN> -e CREDENTIAL_JSON==<BASE64_DATA_OF_CREDENTIAL_JSON_FILE>
```

**Note**: the OMG CLI requires [Docker](https://docs.docker.com/install/) to be installed.

## License
[MIT License](https://github.com/heaptracetechnology/gmail/blob/master/LICENSE).
