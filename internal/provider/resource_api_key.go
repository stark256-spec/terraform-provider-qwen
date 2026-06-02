package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &APIKeyResource{}

type APIKeyResource struct{ client *QwenClient }

type APIKeyResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Status     types.String `tfsdk:"status"`
	PartialKey types.String `tfsdk:"partial_key"`
	SecretKey  types.String `tfsdk:"secret_key"`
	CreatedAt  types.String `tfsdk:"created_at"`
}

func NewAPIKeyResource() resource.Resource { return &APIKeyResource{} }

func (r *APIKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *APIKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":        schema.StringAttribute{Required: true},
			"status":      schema.StringAttribute{Computed: true},
			"partial_key": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"secret_key":  schema.StringAttribute{Computed: true, Sensitive: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"created_at":  schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		},
	}
}

func (r *APIKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil { return }
	client, ok := req.ProviderData.(*QwenClient)
	if !ok { resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("got %T", req.ProviderData)); return }
	r.client = client
}

func (r *APIKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan APIKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() { return }
	k, err := r.client.CreateAPIKey(ctx, plan.Name.ValueString())
	if err != nil { resp.Diagnostics.AddError("Create failed", err.Error()); return }
	plan.ID = types.StringValue(k.ID)
	plan.Status = types.StringValue(k.Status)
	plan.PartialKey = types.StringValue(k.PartialKey)
	plan.CreatedAt = types.StringValue(k.CreatedAt)
	if k.Key != nil { plan.SecretKey = types.StringValue(*k.Key) }
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *APIKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state APIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() { return }
	k, err := r.client.GetAPIKey(ctx, state.ID.ValueString())
	if err != nil { resp.State.RemoveResource(ctx); return }
	state.Name = types.StringValue(k.Name)
	state.Status = types.StringValue(k.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *APIKeyResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "API keys are immutable. Use lifecycle{create_before_destroy=true} to rotate.")
}

func (r *APIKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state APIKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() { return }
	if err := r.client.DeleteAPIKey(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}
