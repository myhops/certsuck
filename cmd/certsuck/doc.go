/*
certsuck is a cli to test and collect certificates from a server that supports tls.

Usage of certsuck:
  -cacerts string
        Extra cacerts to use to verify the server, pem format
  -der-dir string
        Path to write the der files to. Defaults to the current directory (default ".")
  -der-out
        Output der files. [false]
  -der-prefix string
        Prefix for the der files. Defaults to <host name>-
  -host string
        Hostname plus port, e.g. www.example.com:443
  -insecure
        Allow insecure server
  -no-root
        Omit the root cert in pem or der output [false]
  -no-server
        Omit the server cert in pem or der output [false]
  -out
        Show pem output [false]
  -show-opts
        Show the options [false]

-host is the only mandatory flag.

This is the result when used on go.dev:443

 $ certsuck -host go.dev:443
 Verified 0
 0  Subject: CN=go.dev
    Issuer:  CN=WR3,O=Google Trust Services,C=US
 1  Subject: CN=WR3,O=Google Trust Services,C=US
    Issuer:  CN=GTS Root R1,O=Google Trust Services LLC,C=US
 2  Subject: CN=GTS Root R1,O=Google Trust Services LLC,C=US
    Issuer:  CN=GTS Root R1,O=Google Trust Services LLC,C=US
 Verified 1
 0  Subject: CN=go.dev
    Issuer:  CN=WR3,O=Google Trust Services,C=US
 1  Subject: CN=WR3,O=Google Trust Services,C=US
    Issuer:  CN=GTS Root R1,O=Google Trust Services LLC,C=US
 2  Subject: CN=GTS Root R1,O=Google Trust Services LLC,C=US
    Issuer:  CN=GlobalSign Root CA,OU=Root CA,O=GlobalSign nv-sa,C=BE
 3  Subject: CN=GlobalSign Root CA,OU=Root CA,O=GlobalSign nv-sa,C=BE
    Issuer:  CN=GlobalSign Root CA,OU=Root CA,O=GlobalSign nv-sa,C=BE
 Peer
 0  Subject: CN=go.dev
    Issuer:  CN=WR3,O=Google Trust Services,C=US
 1  Subject: CN=WR3,O=Google Trust Services,C=US
    Issuer:  CN=GTS Root R1,O=Google Trust Services LLC,C=US
 2  Subject: CN=GTS Root R1,O=Google Trust Services LLC,C=US
    Issuer:  CN=GlobalSign Root CA,OU=Root CA,O=GlobalSign nv-sa,C=BE
*/
package main

