package routeros

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDhcpLeaseSchemaSupportsMakeStatic(t *testing.T) {
	resource := ResourceDhcpServerLease()
	field, ok := resource.Schema["make_static"]
	if !ok {
		t.Fatal("DHCP lease schema must expose make_static")
	}
	if !field.Optional || field.Type != schema.TypeBool {
		t.Fatal("make_static must be an optional boolean")
	}

	skip := resource.Schema[MetaSkipFields].Default.(string)
	if !strings.Contains(skip, "make_static") {
		t.Fatal("make_static must not be serialized as a RouterOS lease property")
	}
}

func TestMakeStaticTransportMappings(t *testing.T) {
	if got := apiMethodName[crudMakeStatic]; got != "/make-static" {
		t.Fatalf("API make-static command = %q", got)
	}
	if got := restMethodName[crudMakeStatic]; got != "POST" {
		t.Fatalf("REST make-static method = %q", got)
	}
}
