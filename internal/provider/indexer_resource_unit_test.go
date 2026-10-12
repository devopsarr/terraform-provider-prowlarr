package provider

import (
	"context"
	"testing"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/devopsarr/terraform-provider-prowlarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The fields set must be computed-capable: Prowlarr may legitimately return
// normalized values or additional default fields, and a Required-only set
// makes Terraform reject those API responses as inconsistent results.
func TestIndexerFieldsAttributeIsComputedCapable(t *testing.T) {
	t.Parallel()

	resp := &resource.SchemaResponse{}
	(&IndexerResource{}).Schema(context.Background(), resource.SchemaRequest{}, resp)
	require.False(t, resp.Diagnostics.HasError())

	attr, ok := resp.Schema.Attributes["fields"].(schema.SetNestedAttribute)
	require.True(t, ok, "fields attribute should be a set nested attribute")

	assert.False(t, attr.Required, "fields must not be Required-only")
	assert.True(t, attr.Optional, "fields should remain settable from configuration")
	assert.True(t, attr.Computed, "fields must be Computed so API-returned values are accepted")
}

// State must mirror the API response: fields returned only by the API are
// kept, and masked sensitive values are reconciled from the planned value,
// which is the only source that knows the real secret.
func TestIndexerWriteMirrorsAPIResponse(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	planned := []Field{
		{
			Name:           types.StringValue("baseSettings.apiKey"),
			TextValue:      types.StringNull(),
			SensitiveValue: types.StringValue("real-secret"),
			NumberValue:    types.NumberNull(),
			BoolValue:      types.BoolNull(),
			SetValue:       types.SetNull(types.Int64Type),
		},
	}

	planSet, diags := types.SetValueFrom(ctx, IndexerResource{}.getFieldSchema().Type(), planned)
	require.False(t, diags.HasError())

	masked := prowlarr.NewField()
	masked.SetName("baseSettings.apiKey")
	masked.SetValue(helpers.SensitiveValue)

	apiOnly := prowlarr.NewField()
	apiOnly.SetName("baseSettings.queryLimit")
	apiOnly.SetValue("100")

	response := prowlarr.NewIndexerResource()
	response.SetFields([]prowlarr.Field{*masked, *apiOnly})

	indexer := &Indexer{Fields: planSet}
	indexer.write(ctx, response, &diags)
	require.False(t, diags.HasError())

	var state []Field
	require.False(t, indexer.Fields.ElementsAs(ctx, &state, false).HasError())
	require.Len(t, state, 2, "API-only default field must be kept in state")

	byName := make(map[string]Field, len(state))
	for _, f := range state {
		byName[f.Name.ValueString()] = f
	}

	assert.Equal(t, "100", byName["baseSettings.queryLimit"].TextValue.ValueString())
	assert.Equal(t, "real-secret", byName["baseSettings.apiKey"].SensitiveValue.ValueString(),
		"masked API value must be reconciled from the planned sensitive value")
	assert.True(t, byName["baseSettings.apiKey"].TextValue.IsNull())
}
