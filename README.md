# intercert [![Build Status](https://travis-ci.org/evenh/intercert.svg?branch=master)](https://travis-ci.org/evenh/intercert) [![Go Report Card](https://goreportcard.com/badge/github.com/evenh/intercert)](https://goreportcard.com/report/github.com/evenh/intercert) [![codecov](https://codecov.io/gh/evenh/intercert/branch/master/graph/badge.svg)](https://codecov.io/gh/evenh/intercert)

_Brings Let's Encrypt to LAN and other locked down environments._

---
**This is a work in progress (unstable). Contributions are very welcome!**

## How it works

1. A server instance is running somewhere in your network infrastructure, with network access to
your DNS [provider of choice](https://github.com/go-acme/lego/tree/master/providers/dns) and the ACME directory you'll want to use (Let's Encrypt most likely).    
The server is configured with the DNS names you control (e.g. `somecompany.io` and `other.co`).
2. Clients are deployed on the machines where you need the certificates for your applications.
3. Certificates magically appear on the client machine in the directory you've configured.

## Deployment diagram

```
                                                                                        
                                             LAN                                        
  +------------------------------------------------------------------------------------+
  |                                                                                    |
  |                                                                                    |
  |                                                                                    |
  |                     Server 1                                                       |
  |  +--------------------------------------------+                                    |
  |  |                                            |                                    |
  |  | my-db.somecompany.io                       |                                    |
  |  |  app1.somecompany.io   intercert (client)  |                                    |
  |  |  app2.somecompany.io                       |          +-----------------------+ |
  |  +--------------------------------------------+----------|                       | |
  |                                                          |                       | |
  |                     Server N                             |   intercert (server)  | |
  |  +--------------------------------------------+----------|                       | |
  |  |                                            |          +-----------------------+ |
  |  | redis.somecompany.io                       |           /                   |    |
  |  |    intranet.other.co   intercert (client)  |          /                    |    |
  |  |                                            |         /                     |    |
  |  +--------------------------------------------+        /                      |    |
  |                                                       /                       |    |
  +------------------------------------------------------/------------------------|----+
                                      +------------------        +----------------|-+   
                                      |                 |        |                  |   
                                      |   DNS-provider  |        |   ACME provider  |   
                                      |                 |        |                  |   
                                      +-----------------+        +------------------+                                                   
```


## Thanks

A huge thanks to these projects

- [certmagic](https://github.com/mholt/certmagic) - does the hard work for intercert
- [lego](https://github.com/go-acme/lego) - the underpinning library for certmagic, and provides the DNS validation capability 
