# protoc-gen-terraform

protoc plugin that generates Terraform Plugin Framework schema definitions and marshalling/copy functions for golang/protobuf .proto files.

This project is HEAVILY inspired by [Gravitational's protoc-gen-terraform](https://github.com/gravitational/protoc-gen-terraform). However, it is entirely using golang/protobuf and the protogen framework as gogo/protobuf is deprecated.

## Installation

```sh
go install github.com/liamawhite/protoc-gen-terraform
```

## Usage

Running the plugin is pretty typical for a protoc plugin. See the [Makefile](./Makefile) for an example.

### Annotations

| Behavior | Annotation |
| ---- | ---------- |
| `Optional` | Set to true if neither the `REQUIRED` or `OUTPUT_ONLY` [field behavior](https://github.com/googleapis/googleapis/blob/master/google/api/field_behavior.proto#L61) is used. |
| `Required` | Set to true if the `REQUIRED` [field behavior](https://github.com/googleapis/googleapis/blob/master/google/api/field_behavior.proto#L61) is used. |
| `Computed` | Set to true if the `OUTPUT_ONLY` [field behavior](https://github.com/googleapis/googleapis/blob/master/google/api/field_behavior.proto#L61) is used. |

Examples can be found in the [test directory](./test/primary.proto).




