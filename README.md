# Cert Suck

This package offers a simple cli tool that allows you to collect the certificates that the client uses to validate the tls endpoint of a server. 
It can be used to collect the ca chain and create a trust store for e.g. Java applications.

By default it shows the subjects and issuers of the certificates that validated the server.
It support options to show the certificate contents in PEM format, and to write out the DER representations of the certificates.
It also has flags to leave out the certificate of the server and the root certificate, leaving only the intermediate certs if present.

Use the -h option to get the usage information.

If you have any comments or feature requests, please create an issue. 
This program was created quickly and has a lot of room for improvement, I know.

## Usage

Run ```certsuck -h``` to get the usage text.

```bash
$ certsuck -h 
Usage of /certsuck:
  -der-out
        Output der files. [false]
  -der-prefix string
        Prefix for the der files. Defaults to <host name>-
  -host string
        Hostname plus port
  -no-root
        Do not show the root cert in pem output [false]
  -no-server
        Do not show the server cert in pem output [false]
  -out
        Show pem output [false]
  -show-opts
        Show the options [false]
```

### Default use

Run certsuck with the -host option. 
This option expects a hostname plus port number.

The output shows the ca chains for the server.

```
$ certsuck -host jira.belastingdienst.nl:443
Chain 0
  0 Subject: CN=wildcard.belastingdienst.nl,O=ODC Belastingdienst,C=NL
    Issuer:  CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
  1 Subject: CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
    Issuer:  CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
Chain 1
  0 Subject: CN=wildcard.belastingdienst.nl,O=ODC Belastingdienst,C=NL
    Issuer:  CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
  1 Subject: CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
    Issuer:  CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
  2 Subject: CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
    Issuer:  CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
```

By default certsuck shows the names of the certificates that were used to validate the server.
In some cases it shows two chains. 
The longer chain shows the server certificate up and including the root cert.

### Show PEM blocks

One of the main reasons for the existence of this tool is to collect missing certificates and to easily build truststores for Java applications.

To output PEM certificates, use the -out flag in addition to the -host flag.
This will output the certificates in the longest chain.
This includes the server and the root certificate.

```bash
$ certsuck -host jira.belastingdienst.nl:443 -out
Chain 0
  0 Subject: CN=wildcard.belastingdienst.nl,O=ODC Belastingdienst,C=NL
    Issuer:  CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
  1 Subject: CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
    Issuer:  CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
Chain 1
  0 Subject: CN=wildcard.belastingdienst.nl,O=ODC Belastingdienst,C=NL
    Issuer:  CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
  1 Subject: CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
    Issuer:  CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
  2 Subject: CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
    Issuer:  CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
Server certificate
0  Subject: CN=wildcard.belastingdienst.nl,O=ODC Belastingdienst,C=NL
   Issuer:    CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
-----BEGIN CERTIFICATE-----
MIIF+TCCA+GgAwIBAgIUVtUShamvDRKVhvrivhVtMyAZ7PgwDQYJKoZIhvcNAQEL
BQAwTDELMAkGA1UEBhMCTkwxHDAaBgNVBAoME09EQyBCZWxhc3RpbmdkaWVuc3Qx
HzAdBgNVBAMMFkluZnJhc3RydWN0dXVyIENBIC0gRzMwHhcNMjQwMTMwMTIzNTUz
WhcNMjYwNTAzMTIzNTUyWjBRMQswCQYDVQQGEwJOTDEcMBoGA1UECgwTT0RDIEJl

... Rest of the output is omitted.
```

You can use the ```-no-root``` and ```-no-server``` option to omit the root and the server certificate and keep the intermediate certs.
This set of certificates is usually sufficient for a trust store.
The output can be saved to a file using redirection.

### Write DER files

To create a truststore.jks file, you need a DER representation of the certificates. 
Use the ```-out-der``` option write each certificate in a separate file in the current directory.
The names of the files are the name of the host followed by the index in the chain, followed by **.der**.
Use the ```-der-prefix``` to use a different prefix for the filename.
Use the ```-der-dir``` option to write the files to a different directory.

```bash
$ certsuck -host jira.belastingdienst.nl:443 -der-out -der-dir der-dir -der-prefix der-prefix- 
Chain 0
  0 Subject: CN=wildcard.belastingdienst.nl,O=ODC Belastingdienst,C=NL
    Issuer:  CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
  1 Subject: CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
    Issuer:  CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
Chain 1
  0 Subject: CN=wildcard.belastingdienst.nl,O=ODC Belastingdienst,C=NL
    Issuer:  CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
  1 Subject: CN=Infrastructuur CA - G3,O=ODC Belastingdienst,C=NL
    Issuer:  CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
  2 Subject: CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL
    Issuer:  CN=ODC Belastingdienst Root CA - G1,O=ODC Belastingdienst,C=NL

$ ls der-dir -l
total 12
-rwxr-xr-x 1 zandp06 domain users 1533 jul  3 17:05 jira.belastingdienst.nl-00.der
-rwxr-xr-x 1 zandp06 domain users 1649 jul  3 17:05 jira.belastingdienst.nl-01.der
-rwxr-xr-x 1 zandp06 domain users 1400 jul  3 17:05 jira.belastingdienst.nl-02.der
```

### Other options

* Use ```-h``` to show the usage.
* Use ```-show-opts``` to show the used options.

## Installation

If you have go installed, run ```go install github.com/myhops/certsuck/cmd/certsuck@latest```.

