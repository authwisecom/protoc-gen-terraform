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

package test

import (
	"context"
	"testing"

	types "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestSchema(t *testing.T) {
	schema, diags := GenSchemaTest(context.Background())

	// Check no errors
	require.False(t, diags.HasError())

	t.Run("Field Annotations", func(*testing.T) {
		require.True(t, schema.Attributes["required"].Required)
		require.True(t, schema.Attributes["str"].Optional)
	})

	t.Run("Field injection", func(*testing.T) {
		require.True(t, schema.Attributes["inject_computed"].Computed)
		require.True(t, schema.Attributes["inject_required"].Required)
		require.True(t, schema.Attributes["inject_optional"].Optional)
	})

	t.Run("Primitive types", func(*testing.T) {
		require.Equal(t, types.BoolType, schema.Attributes["bool"].Type)
		require.Equal(t, types.Float64Type, schema.Attributes["double"].Type)
		require.Equal(t, types.Float64Type, schema.Attributes["float"].Type)
		require.Equal(t, types.Int64Type, schema.Attributes["int32"].Type)
		require.Equal(t, types.Int64Type, schema.Attributes["int64"].Type)
		require.Equal(t, types.StringType, schema.Attributes["str"].Type)
		require.Equal(t, types.StringType, schema.Attributes["bytes"].Type)
	})

	t.Run("List with primitive type", func(*testing.T) {
		require.Equal(t, types.ListType{ElemType: types.StringType}, schema.Attributes["string_list"].Type)
		require.Nil(t, schema.Attributes["string_list"].Attributes)
	})

	t.Run("Map with primitive type", func(*testing.T) {
		require.Equal(t, types.MapType{ElemType: types.StringType}, schema.Attributes["map"].Type)
		require.Nil(t, schema.Attributes["map"].Attributes)
	})

	t.Run("Nested message", func(*testing.T) {
		// Single
		require.Equal(t, types.StringType, schema.Attributes["nested"].Attributes.GetAttributes()["str"].GetType())

		// List
		require.Nil(t, schema.Attributes["nested_list"].Type)
		require.Equal(t, types.StringType, schema.Attributes["nested_list"].Attributes.GetAttributes()["str"].GetType())

		// Map
		require.Nil(t, schema.Attributes["nested_map"].Type)
		require.Equal(t, types.StringType, schema.Attributes["nested_map"].Attributes.GetAttributes()["str"].GetType())
	})

	t.Run("Enum", func(*testing.T) {
		require.Equal(t, types.Int64Type, schema.Attributes["mode"].Type)
	})

	t.Run("OneOfs", func(*testing.T) {
		require.Nil(t, schema.Attributes["branch1"].Type)
		require.Equal(t, types.StringType, schema.Attributes["branch1"].Attributes.GetAttributes()["str"].GetType())
		require.Nil(t, schema.Attributes["branch2"].Type)
		require.Equal(t, types.Int64Type, schema.Attributes["branch2"].Attributes.GetAttributes()["int32"].GetType())
		require.Nil(t, schema.Attributes["branch3"].Attributes)
		require.Equal(t, types.StringType, schema.Attributes["branch3"].Type)
	})
}

func TestSchemaMultipleFiles(t *testing.T) {
	schema, diags := GenSchemaTest2(context.Background())
	require.False(t, diags.HasError())
	require.Equal(t, types.StringType, schema.Attributes["str"].Type)
}
