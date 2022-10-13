// Copyright 2022 Liam White
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generate

import (
	"regexp"
	"strings"

	j "github.com/dave/jennifer/jen"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func Scheme(f *j.File, m *protogen.Message) {
	id := "GenSchema" + m.GoIdent.GoName
	f.Commentf("// %v returns tfsdk.Schema definition for %v\n", id, m.GoIdent.GoName).
		Func().
		Id(id).
		Params(j.Id("ctx").Qual("context", "Context")).
		Params(j.Qual(SDK, "Schema"), j.Qual(Diag, "Diagnostics")).
		Block(j.Return(
			j.Qual(SDK, "Schema").Values(j.Dict{
				j.Id("Attributes"): j.Map(j.String()).Qual(SDK, "Attribute").Values(
					fieldsDictSchema(m, true),
				),
			}),
			j.Nil(),
		))
}

func fieldsDictSchema(m *protogen.Message, injectId bool) j.Dict {
	d := j.Dict{}

	for _, f := range m.Fields {
		d[j.Lit(snakeCase(f.GoName))] = field(f)
	}

	if injectId {
		d[j.Lit("id")] = j.Values(j.Dict{
			j.Id("Type"):     j.Qual(Types, "StringType"),
			j.Id("Computed"): j.Lit(true),
		})
	}

	return d
}

func field(f *protogen.Field) j.Code {
	d := j.Dict{
		j.Id("Description"): j.Lit(trimComments(f.Comments.Leading)),
		j.Id("Type"):        schemaType(f.Desc), // nils are automatically omitted
		j.Id("Attributes"):  attributes(f),
	}

	// Handle field behavior annotations
	opts := f.Desc.Options().(*descriptorpb.FieldOptions)
	optional := true
	for _, b := range proto.GetExtension(opts, annotations.E_FieldBehavior).([]annotations.FieldBehavior) {
		switch b {
		case annotations.FieldBehavior_REQUIRED:
			d[j.Id("Required")] = j.Lit(true)
			optional = false
		case annotations.FieldBehavior_OUTPUT_ONLY:
			d[j.Id("Computed")] = j.Lit(true)
			optional = false
		}
	}
	// If required or computed is not set, default to optional
	if optional {
		d[j.Id("Optional")] = j.Lit(true)
	}

	return j.Values(d)
}

var primitiveTypeMap = map[protoreflect.Kind]*j.Statement{
	protoreflect.StringKind: j.Qual(Types, "StringType"),
	protoreflect.BytesKind:  j.Qual(Types, "StringType"),
	protoreflect.Int32Kind:  j.Qual(Types, "Int64Type"),
	protoreflect.Int64Kind:  j.Qual(Types, "Int64Type"),
	protoreflect.EnumKind:   j.Qual(Types, "Int64Type"),
	protoreflect.FloatKind:  j.Qual(Types, "Float64Type"),
	protoreflect.DoubleKind: j.Qual(Types, "Float64Type"),
	protoreflect.BoolKind:   j.Qual(Types, "BoolType"),
}

func schemaType(d protoreflect.FieldDescriptor) *j.Statement {
	if d.IsList() {
		// If the type isnt a primitive then type is nil, we use attributes instead.
		if _, ok := primitiveTypeMap[d.Kind()]; !ok {
			return nil
		}
		return j.Qual(Types, "ListType").Values(j.Dict{
			j.Id("ElemType"): primitiveTypeMap[d.Kind()],
		})
	}
	if d.IsMap() {
		// If the type isnt a primitive then type is nil, we use attributes instead.
		if _, ok := primitiveTypeMap[d.MapValue().Kind()]; !ok {
			return nil
		}
		return j.Qual(Types, "MapType").Values(j.Dict{
			j.Id("ElemType"): primitiveTypeMap[d.MapValue().Kind()],
		})
	}
	return primitiveTypeMap[d.Kind()]
}

func attributes(f *protogen.Field) *j.Statement {
	// If message is not nil it can't be a primitive type (string, bool, etc.).
	if f.Message != nil {
		if f.Desc.IsList() {
			return xNestAttributes("List", f.Message)
		}
		if f.Desc.IsMap() {
			// If the map has a primitive value we use type, not attributes.
			if _, ok := primitiveTypeMap[f.Desc.MapValue().Kind()]; ok {
				return nil
			}
			// Not sure how safe the assumption that fields[1] is always value and not key ¯\_(ツ)_/¯.
			return xNestAttributes("Map", f.Message.Fields[1].Message)
		}
		// If we've got this far is must be single nested
		return xNestAttributes("Single", f.Message)

	}
	return nil
}
func xNestAttributes(typ string, m *protogen.Message) *j.Statement {
	return j.Qual(SDK, typ+"NestedAttributes").Params(
		j.Map(j.String()).Qual(SDK, "Attribute").Values(fieldsDictSchema(m, false)),
	)
}

func trimComments(c protogen.Comments) string {
	return strings.TrimSpace(strings.TrimPrefix(c.String(), "// "))
}

func snakeCase(s string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
