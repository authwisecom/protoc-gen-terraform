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

	// Check field annotations
	require.True(t, schema.Attributes["computed"].Computed)
	require.True(t, schema.Attributes["required"].Required)
	require.True(t, schema.Attributes["str"].Optional)

	// Check ID field was injected in top level messages only
	require.Equal(t, types.StringType, schema.Attributes["id"].Type)
	require.Nil(t, schema.Attributes["nested"].Attributes.GetAttributes()["Id"])

	// Check each primitive type
	require.Equal(t, types.BoolType, schema.Attributes["bool"].Type)
	require.Equal(t, types.Float64Type, schema.Attributes["double"].Type)
	require.Equal(t, types.Float64Type, schema.Attributes["float"].Type)
	require.Equal(t, types.Int64Type, schema.Attributes["int32"].Type)
	require.Equal(t, types.Int64Type, schema.Attributes["int64"].Type)
	require.Equal(t, types.StringType, schema.Attributes["str"].Type)
	require.Equal(t, types.StringType, schema.Attributes["bytes"].Type)

	// Check List with primitive type
	require.Equal(t, types.ListType{ElemType: types.StringType}, schema.Attributes["string_list"].Type)
	require.Nil(t, schema.Attributes["string_list"].Attributes)

	// Check Map with primitive type
	require.Equal(t, types.MapType{ElemType: types.StringType}, schema.Attributes["map"].Type)
	require.Nil(t, schema.Attributes["map"].Attributes)

	// Check Nested
	require.Equal(t, types.StringType, schema.Attributes["nested"].Attributes.GetAttributes()["str"].GetType())

	// Check NestedList
	require.Nil(t, schema.Attributes["nested_list"].Type)
	require.Equal(t, types.StringType, schema.Attributes["nested_list"].Attributes.GetAttributes()["str"].GetType())

	// Check NestedMap
	require.Nil(t, schema.Attributes["nested_map"].Type)
	require.Equal(t, types.StringType, schema.Attributes["nested_map"].Attributes.GetAttributes()["str"].GetType())

	// Check Enum
	require.Equal(t, types.Int64Type, schema.Attributes["mode"].Type)

	// Check OneOf
	require.Nil(t, schema.Attributes["branch1"].Type)
	require.Equal(t, types.StringType, schema.Attributes["branch1"].Attributes.GetAttributes()["str"].GetType())
	require.Nil(t, schema.Attributes["branch2"].Type)
	require.Equal(t, types.Int64Type, schema.Attributes["branch2"].Attributes.GetAttributes()["int32"].GetType())
	require.Nil(t, schema.Attributes["branch3"].Attributes)
	require.Equal(t, types.StringType, schema.Attributes["branch3"].Type)
}

func TestSchemaMultipleFiles(t *testing.T) {
	schema, diags := GenSchemaTest2(context.Background())
	require.False(t, diags.HasError())
	require.Equal(t, types.StringType, schema.Attributes["str"].Type)
}
