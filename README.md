# Go-WebAuth

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