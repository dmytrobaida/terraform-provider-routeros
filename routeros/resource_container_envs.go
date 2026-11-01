package routeros

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// https://help.mikrotik.com/docs/display/ROS/Container#Container-Addenvironmentvariablesandmounts(optional)
func ResourceContainerEnvs() *schema.Resource {
	resSchema := map[string]*schema.Schema{
		MetaResourcePath: PropResourcePath("/container/envs"),
		MetaId:           PropId(Id),

		"list": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Environment variables list (RouterOS 7.23 and newer).",
		},
		"key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Key of the environment variable.",
		},
		"value": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Value of the environment variable.",
		},
	}

	return &schema.Resource{
		CreateContext: DefaultCreate(resSchema),
		ReadContext:   DefaultRead(resSchema),
		UpdateContext: DefaultUpdate(resSchema),
		DeleteContext: DefaultDelete(resSchema),

		Importer: &schema.ResourceImporter{
			StateContext: ImportStateCustomContext(resSchema),
		},

		Schema: resSchema,
	}
}
