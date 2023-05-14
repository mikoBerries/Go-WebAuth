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
●Base 64 are not used to encryption for stuff but base64 have a plus because it's printable / readble. so we can throw it everywhere.

●HMAC: is a specific type of message authentication code (MAC)

SHA
----
●SHA-256 is used in some of the most popular authentication and encryption protocols, including SSL, TLS, IPsec, SSH, and PGP. In Unix and Linux, SHA-256 is used for secure password hashing.