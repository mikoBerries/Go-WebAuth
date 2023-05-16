# Go-WebAuth
authentication vs authorization
--------------------
● authentication (AuthN) : verify that someone or something is who or what they claim to be.
● authorization (AuthZ): Authorization is the security process that determines a user or service's level of access

ROTATING PRIVATE KEY
--------------------
● you will not want to use same private key same all the time then we need to change key just like social media password .
● The idea is to make new combination since we don't know the older has been cracked or not.

Signing
-------
● Symmetric Key (signing using 1 key got for data that only you will see ex database ,logging)
   ● HMAC 
● Asymmetric Key using 2 key private and public key
    RSA,ECDSA private key to sign(encrypt)/public key to verify(decrypt)
● other signing method that use Asymmetric Key to send Symmetric Key to each other, then every body just using 1 key to comunicated privately (Ecnrypt & Decrypt)


BASE64
----------------
● Base 64 are not used to encryption for stuff but base64 have a plus because it's printable / readble. so we can throw it everywhere.

● HMAC: is a specific type of message authentication code (MAC)

SHA
----
● SHA-256 is used in some of the most popular authentication and encryption protocols, including SSL, TLS, IPsec, SSH, and PGP. In Unix and Linux, SHA-256 is used for secure password hashing.


JWT (Json Web Token)
-----------------------

● A stateless authentication mechanism is one where the server does not need to maintain any state about the user's session. This means that the server does not need to store any information about the user, such as their username, password, or session ID.

● JWT is a stateless authentication mechanism because the token itself contains all of the information that the server needs to authenticate the user. The token is signed by the server, so the client can be confident that it is authentic.

● JThe common information that the server needs to authenticate the user that store in JWT is:

○ User ID: This is the unique identifier for the user.
○ Username: This is the username of the user.
○ Email: This is the email address of the user.
○ Role: This is the role of the user.
○ Expiration date: This is the date and time after which the JWT will no longer be valid.

The server can use this information to authenticate the user and determine whether they have the appropriate permissions to access the requested resource.

It is important to note that the JWT should not contain any sensitive information, such as the user's password. This is because the JWT is sent to the client in the clear, and it could be intercepted by a malicious actor.

●  JWT can be stored in cookie if it's web client or local storage if it's from mobile.

●  For best practice store at jwt sessionId that reference to DB data so we can changes permission ASAP with update DB. 

OAuth2 (https://aaronparecki.com/oauth-2-simplified/#authorization)
-------------------------------------------



Graphql
--------
http://spec.graphql.org/June2018/#sec-Overview

There are three types of operations that GraphQL models:

query – a read‐only fetch. --> in oauth login we just using read only data
mutation – a write followed by a fetch.
subscription – a long‐lived request that fetches data in response to source events.


few useful library
https://github.com/bashtian/jsonutils  for json
https://transform.tools/json-to-go