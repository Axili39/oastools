OAS Tool generator demo
=======================

```
../oatoolgen -f demo.yaml
tree .

.
├── cmd
│   └── demoTool.go
├── demo
│   ├── demo.pb.go
│   └── demo.proto
├── demo.yaml
└── README.md

go build cmd/demoTool.go
./demoTool -h

Usage of ./demoTool:
  -f string
        input file .json/.yaml/.bin
  -g    generate empty file
  -o string
        json|yaml|bin (default "bin")
  -r string
```