.PHONY: build clean test

clean:
	rm -f test/*.pb.go
	rm -f test/*_terraform.go


build:
	go install github.com/liamawhite/protoc-gen-terraform
	protoc -Iextensions/google/api -I. --go_out=. --go_opt=paths=source_relative --terraform_out=. --terraform_opt=paths=source_relative test/primary.proto test/secondary.proto

test: clean build
	go test ./...  
