# Feature ideas

- live-reloading development server
- production webserver using autocert to get certs
- static site generator
  - call using -gen flag to generate the site
  - prefer embedding css and images 
- grpc integration, with live-reload on proto file changes
- gopherjs integration
- use internal package that finds the pwd's current go "project" but searching upwards for a git directory
- should be able to work using either `go run` or building a binary
- has template for most of SaaS common features. home page, marketting page
  


# Func ideas

- include builders to make common templated sites
  - Define "product" once, generate sales site, documentation site, support site all with sane defaults.  Customize with modifier functions.
  - Open source project renders home page, links to documentation, etc.  Use github jekyl site as inspiration
- handlers can return a specific error type to halt additional processing


# Plugin ideas
- callback-based plugin system? middleware-based plugin system? both?
- plugin package has init function to call registration funcs on sites package
- 

# Questions

- Can we use go-bindata to embed built files?
- What can we leverage from hugo?
- What benefit do we get from defining a grpc backend?
  - Is "sites", actually "products" or "services"?
    - 
- should the whole thing be plugin based?
- how does a product get configured to communicate with other grpc services, such as an "auth" service?
  - perhaps a plugin that declared handlers that handles omniauth-style callbacks?
  - 

# Example packages

First thoughts:
```golang
package main

import (
  "github.com/nullstyle/go/produce
)

// declare the new product.  Give it a name for when it appears in dashboards or admin pages.
var product = produce.New("basic-site")

func init() {
    // raw handlers:
    //   - when server, serve directly
    //   - when building copy directly
    //   - when developing, serve directly
    // Under means trigger on all children, recursively
    product.Raw("/public").At("/")

    // any additinal args are handlers
    // Grpc handler:
    //   - when server, runs the grpc server
    //   - when building, adds context to build environment for production api server endpoint
    //   - when developing, live-reloading grpc server that restarts when the proto file changes or the server implementation changes
    site.Grpc("basic-api")
}

```