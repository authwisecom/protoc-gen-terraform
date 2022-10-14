# Copyright 2022 Liam White
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: build clean test

clean:
	rm -f test/*.pb.go
	rm -f test/*_terraform.go

build:
	go install github.com/liamawhite/protoc-gen-terraform
	protoc -Iextensions/google/api -Iextensions/google/protobuf -I. --go_out=. --go_opt=paths=source_relative --terraform_out=. --terraform_opt=paths=source_relative  --terraform_opt=loglevel=0 test/primary.proto test/secondary.proto

test: clean build
	go test ./...  

format:
	./ci/format
	./ci/dirtystate