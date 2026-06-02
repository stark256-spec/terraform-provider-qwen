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

var _ resource.Resource = &WorkspaceResource{}

type WorkspaceResource struct{ client *QwenClient }

type WorkspaceResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func NewWorkspaceResource() resource.Resource { return &WorkspaceResource{} }

func (r *WorkspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (r *WorkspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":       schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"status":     schema.StringAttribute{Computed: true},
			"created_at": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		},
	}
}

func (r *WorkspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil { return }
	client, ok := req.ProviderData.(*QwenClient)
	if !ok { resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("got %T", req.ProviderData)); return }
	r.client = client
}

func (r *WorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkspaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() { return }
	w, err := r.client.CreateWorkspace(ctx, plan.Name.ValueString())
	if err != nil { resp.Diagnostics.AddError("Create failed", err.Error()); return }
	plan.ID = types.StringValue(w.ID)
	plan.Status = types.StringValue(w.Status)
	plan.CreatedAt = types.StringValue(w.CreatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkspaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() { return }
	w, err := r.client.GetWorkspace(ctx, state.ID.ValueString())
	if err != nil { resp.Diagnostics.AddError("Read failed", err.Error()); return }
	state.Name = types.StringValue(w.Name)
	state.Status = types.StringValue(w.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkspaceResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Workspace names are immutable. Replace to rename.")
}

func (r *WorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkspaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() { return }
	if err := r.client.DeleteWorkspace(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}
