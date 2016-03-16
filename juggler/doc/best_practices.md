# Best practices

It's still way too early to dare make any best practice recommendations, except for one.

## All RPC and Pub-Sub calls must go through the server

The core of the application's security is handled by the juggler server's middleware handlers (authentication, authorization, rate limiting, URI whitelisting, etc.). If an internal component (e.g. a worker process that listens on a message queue) needs to emit events or call RPC functions, it should connect to the juggler server as a client and use that connection to make requests.

That is to say, it should *never* connect directly to redis and use redis commands to make requests (how juggler uses redis is an implementation detail), and it should *never* directly use a broker to make such requests either. This enforces loose coupling to redis (only accessed indirectly via the server) and ensures the handlers are always executed as designed.
