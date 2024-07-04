/*
certsuck collects certificates from the given server.
  $ certsuck -h
  Usage of certsuck:
  -der-dir string
    Path to write the der files to. Defaults to the current directory (default ".")
  -der-out
    Output der files. The names of the files is <host>-0x.der [false]
  -der-prefix string
    Prefix for the der files. Defaults to <host name>-
  -host string
    Hostname plus port
  -no-root
    Omit the root cert in pem or der output [false]
  -no-server
    Omit the server cert in pem or der output [false]
  -out
    Show pem output [false]
  -show-opts
    Show the options [false]	
*/
package main

