### Easy TLS

#### A Simple web server whith automatic TLS config.  

Using LetsEncrypt generates a server certificate using the http domain validation
so on startup a new key and certificate are created and cached in the certs directory.  

Uses a single Environment variable `EASYTLS_DOMAIN` to define the domain for which the certificate is to be issued.  
This domain must match the domain of the host server in order to validate ownership of that domain.  

If environment variable is not set, will work in 'developer' mode, where a simple, self signed certificate is generated instead,.

Docker image can be found:  `https://hub.docker.com/repository/docker/eurospoofer/easytls`