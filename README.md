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
oa2proto -f FILE [-node component1 ... -node componenentn] [-p package] [-option option1 ... -option optionn] [-o FILE.proto] [-build path] 
