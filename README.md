Golang Open Api Specs tools
============================

Golang package providing OpenApi class model.

Tools:
======
* **oatree**: Dump model as Simple Tree,
* **objtoolgen**: Generate a tool for managing yaml, json or binary encoded files specified by a Open Api Schema.
* **oa2proto**: convert OpenApi spec into protobuf .proto spec file.

Install
-------
```
go get github.com/Axili39/oastools/cmd/oatree
or
go get github.com/Axili39/oastools/cmd/objtoolgen
or
go get github.com/Axili39/oastools/cmd/oa2proto

go get github.com/Axili39/oastools/
```

Usage
-----
oa2proto -f FILE [-node component1 ... -node componenentn] [-p package] [-add-enum-prefix] [-rename-package FILE:PACKAGE ...] [-option option1 ... -option optionn] [-o FILE.proto] [-build path] 

Usage of oa2proto:
  -add-enum-prefix
        Auto add prefix on Enums
  -build string
        build with protoc
  -f string
        yaml file to parse
  -node value
        select component (multi)
  -o string
        output file
  -option value
        add directive option in .proto file (multi)
  -p string
        package name eg: foo.bar
  -rename-package value
        rename package imports
  -v    show version

  Notes:
  ======
  External References
  -------------------
  consider child.yaml :
  ```yaml
  openapi: "3.0.0"
info:
  version: 1.0.0
  title: Lux API
  license:
    name: MIT
servers:
  - url: http://github.com/Axili39/lux
paths:
components:
  schemas:
    bar:
      type: object
      properties:
        member1:
          type: string
        member2:
          description: |
            ligne 1
            ligne 2
          type: integer
  ```
    
and root.yaml which has a reference to child.yaml.
  ```yaml
openapi: "3.0.0"
info:
  version: 1.0.0
  title: Lux API
  license:
    name: MIT
  x-package: root
servers:
  - url: http://github.com/Axili39/lux
paths:
components:
  schemas:
    foo:
      type: object
      properties:
        member1:
          type: string
        member2:
          description: "External object"
          $ref: "child.yaml#/components/schemas/bar"
  ```

We can generate 2 différents .proto files which one import the other :
```sh
oa2proto -f root.yaml -rename-package child.yaml:child -p root -o root.proto -option "go_package=\"gen/root\""
oa2proto -f child.yaml -p child -o child.proto -option="go_package=\"gen/child\""
```

output for child.proto :
```protobuf
syntax = "proto3";
package  child ;
option  go_package="gen/child" ;
/* Type :  */
message bar {
        string member1 = 1; /*  */
        int32 member2 = 2; /* ligne 1
ligne 2
 */
}
```
output for root.proto :
```protobuf
syntax = "proto3";
package  root ;
option  go_package="gen/root" ;
import "child.yaml";
/* Type :  */
message foo {
        string member1 = 1; /* Simple string */
        child.bar member2 = 2; /* External object in child.yaml */
}
```

to compile, use 2 separated commands, example for go_out :
```sh
protoc  -I./ --go_out=gen ./root.proto
protoc  -I./ --go_out=gen ./child.proto
```