Reproducer for unexpected EOFs with "httputil: ReverseProxy read error during body copy: read tcp 127.0.0.1:41980->127.0.0.1:8080: use of closed network connection"
logs in queue-proxy.

* _Set GOPATH_
* _Install go (Makefile expect ${GOPATH}/bin/ko to exist)_
* `make apply-receiver-ksvc`
* _wait a bit_
* `make apply-sender-ksvc`
* Watch logs for the error:
  * `kubectl logs -n sender-ksvc -l serving.knative.dev/service=sender-ksvc -c user-container -f`
  * `kubectl logs -n receiver-ksvc -l serving.knative.dev/service=receiver-ksvc -c queue-proxy -f | grep httputil`
* After a while, notice the errors (maybe)
  * ```httputil: ReverseProxy read error during body copy: read tcp 127.0.0.1:60188->127.0.0.1:8080: use of closed network connection``` in receiver queue-proxy logs
  * ```error reading body: unexpected EOF``` in sender logs
  
There are also `make apply-receiver-k8s` and `make apply-sender-k8s` variants which use plain k8s Deployments and don't manifest the errors.