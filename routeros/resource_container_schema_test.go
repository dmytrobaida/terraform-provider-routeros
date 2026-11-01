package routeros

import "testing"

func TestContainerSchemaSupportsRouterOS723EnvironmentLists(t *testing.T) {
	fields := ResourceContainer().Schema
	modern := fields["envlists"]
	if modern == nil || !modern.Optional || modern.ForceNew || modern.Elem == nil {
		t.Fatal("envlists must be an optional, updatable string set")
	}
	if _, exists := fields["envlist"]; exists {
		t.Fatal("container must omit legacy envlist alias to prevent perpetual drift")
	}
	if !ResourceContainerEnvs().Schema["value"].Sensitive {
		t.Fatal("environment values must be sensitive")
	}
	envs := ResourceContainerEnvs().Schema
	if list := envs["list"]; list == nil || !list.Required || list.Optional {
		t.Fatal("environment resource must require the modern list field")
	}
	if _, exists := envs["name"]; exists {
		t.Fatal("environment resource must omit legacy name to prevent Pulumi auto-naming")
	}
}

func TestContainerSchemaSupportsRouterOS723RestartFields(t *testing.T) {
	schema := ResourceContainer().Schema
	for _, name := range []string{"restart_policy", "restart_interval"} {
		field, ok := schema[name]
		if !ok || !field.Optional || field.ForceNew {
			t.Fatalf("%s must be an optional, updatable field", name)
		}
		if len(field.ConflictsWith) != 1 || field.ConflictsWith[0] != "auto_restart_interval" {
			t.Fatalf("%s must reject the legacy restart setting", name)
		}
	}
	validate := schema["restart_policy"].ValidateFunc
	for _, value := range []string{"no", "on-failure", "always"} {
		if _, errors := validate(value, "restart_policy"); len(errors) != 0 {
			t.Fatalf("valid restart policy %q rejected: %v", value, errors)
		}
	}
	if _, errors := validate("invalid", "restart_policy"); len(errors) == 0 {
		t.Fatal("invalid restart policy accepted")
	}
	if schema["restart_interval"].DiffSuppressFunc == nil {
		t.Fatal("restart interval must normalize RouterOS time values")
	}
	if len(schema["auto_restart_interval"].ConflictsWith) != 2 {
		t.Fatal("legacy restart interval must reject both modern restart settings")
	}
}

func TestContainerSchemaUsesRouterOS721MountFields(t *testing.T) {
	schema := ResourceContainer().Schema

	if _, ok := schema["mountlists"]; !ok {
		t.Fatal("container schema must expose mountlists")
	}
	if _, ok := schema["mounts"]; ok {
		t.Fatal("container schema must not expose obsolete mounts")
	}
}

func TestContainerMountSchemaUsesRouterOS721ListField(t *testing.T) {
	schema := ResourceContainerMounts().Schema

	if _, ok := schema["list"]; !ok {
		t.Fatal("container mount schema must expose list")
	}
	if _, ok := schema["name"]; ok {
		t.Fatal("container mount schema must not expose obsolete name")
	}
}

func TestContainerDefinesAllLifecycleTimeouts(t *testing.T) {
	timeouts := ResourceContainer().Timeouts
	if timeouts == nil || timeouts.Create == nil || timeouts.Update == nil || timeouts.Delete == nil {
		t.Fatal("container resource must define create, update, and delete timeouts")
	}
}

func TestContainerPullAllowsTransientUnknownState(t *testing.T) {
	for _, state := range containerPullPendingStates {
		if state == "unknown" {
			return
		}
	}
	t.Fatal("container pull must tolerate RouterOS returning no lifecycle flag")
}

func TestContainerStateSupportsRouterOS721Flags(t *testing.T) {
	tests := []struct {
		name string
		item MikrotikItem
		want string
	}{
		{name: "legacy status", item: MikrotikItem{"status": "running"}, want: "running"},
		{name: "running flag", item: MikrotikItem{"running": "true"}, want: "running"},
		{name: "extracting flag", item: MikrotikItem{"extracting": "true"}, want: "extracting"},
		{name: "stopped flag", item: MikrotikItem{"stopped": "true"}, want: "stopped"},
		{name: "missing state", item: MikrotikItem{}, want: "unknown"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := containerState(test.item); got != test.want {
				t.Fatalf("containerState() = %q, want %q", got, test.want)
			}
		})
	}
}
