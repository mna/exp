# juggler

TODO :
* benchmarking tool (docker-based?)
* add Close URI to server, test it with client
* redis cluster support and tests
* metrics (expvar?)
* move juggler client to its own package?
* move the msg.Exp pseudo-message to the client package?

## Rationale behind design decisions

This section documents the rationale for some design decisions. Those decisions may be revisited and changed later on, backed by empirical evidence.

### No AUTH mechanism / message

There are many different authentication mechanisms, each with its own advantages and trade-offs, and there are many usage contexts, each with their own security requirements. Many authentication mechanisms can be built in userland using the existing CALL/RES messages and a set of conventions. For example, similar to the cookie with session ID:

* CALL com.example.auth {"username": "x", "password": "y"}
* RES {"authenticated": true, "roles": ["admin"], "token": "xyz"}

Then for URIs that require authentication, the middleware `Handler` can check for a valid token (exists and is not expired), and either proceed with the authenticated call or return:

* ERR {"code": 403, "message": "user is not authenticated"}

A downside is that it makes it harder to connect different juggler clients and servers, since authentication is implementation/application-specific. This is ok given the goal of juggler, which is to build web applications, not for interoperability between heterogeneous systems.

### No CLOSE mechanism / message

Closing a connection can be done cleanly by sending a `close` websocket message from either peer, or it can be done abruptly by closing the socket. Given the nature of network-based communications, each peer must be prepared to handle abruptly closed sockets anyway.

That being said, although clients can send `close` websocket frames, there is no built-in support for server-based close. This can be easily done via the existing CALL/RES messages and a set of conventions. For example:

* The client sends CALL com.example.close
* A middleware handler intercepts the CALL and emits the `close` websocket frame, then closes the connection.

A downside is that it makes it harder to cleanly disconnect from different juggler clients and servers, since the close mechanism is implementation/application-specific. This is ok given the goal of juggler, which is to build web applications, not for interoperability between heterogeneous systems.

### No whitelisting/dynamic discovery of URIs/channels

If a client makes requests to a URI that no callee listens on, it can end up polluting the redis database and eventually use all the redis server's memory (even with a conservative `redisbroker.Broker.CallCap` configuration, since the client can make calls on any number of URIs).

A way to prevent that would be to whitelist the allowed URIs to call. The same is true for the pub-sub channels, although it is less likely to cause dramatic issues (redis presumably quickly drops events that have no subscribers). And a reflection tool that could query which URIs (and to a lesser extent, channels) are available could be useful.

In order to keep the juggler core small and focused, this feature is not provided. However, it can be achieved relatively easily (for URIs at least) in a middleware handler with the help of the callees:

* when a callee starts listening to a URI, it creates a key for that URI, with an expiration (e.g. SET "callee.uri.{uri}" 1 PX 30000).
* the callee starts a goroutine that resets the TTL of the key at a regular interval (less than the key's TTL) as long as it is still listening for that URI (a heartbeat).
* a middleware handler on the server intercepts the CALL messages and checks if the key for that URI exists, returning an ERR if it doesn't.
* a meta-callee listens for a discovery URI (e.g. com.example.ListURIs) and queries the existing keys (e.g. using SCAN), returning the live URIs.

In the worst case, an inactive URI would be reported as "live" for the TTL of the key (if the callee dies just after having set the expiration of the key). Multiple callees can listen for the same URI without problem, and it would stop being reported as "live" only when the last callee stops listening. The dynamic discovery is somewhat more complex in a redis-cluster (it would require a SCAN on each node, see https://github.com/antirez/redis-doc/issues/286).

The downside is that the dynamic reporting and whitelisting of calls only works if the calls go through the server's middleware handler, so internal components that make calls directly using a CallerBroker interface are not subject to those checks, but that's ok because internal components are presumed safe.

