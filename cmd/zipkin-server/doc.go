/*Package main Zipkin API

Zipkin's Query api is rooted at `api/v1`, on a host that by default listens
on port 9411. It primarily serves zipkin-web, although it includes a POST
endpoint that can receive spans.



    Schemes:
      http
      https
    Host: localhost:9411
    BasePath: /
    Version: 1.0.0

    Consumes:
    - application/x-thrift

    - application/json


    Produces:
    - application/json


swagger:meta
*/
package main
