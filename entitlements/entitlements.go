package entitlements

type Entitlement struct {
	Name string `json:"name"`
	// map[roleName][]permissions
	AddedPermissions map[string][]string `json:"permissions"`
}

var entitlements map[string]map[string][]string

func AddEntitlement(e Entitlement) {
	if entitlements == nil {
		entitlements = map[string]map[string][]string{}
	}
	entitlements[e.Name] = e.AddedPermissions
}

func GetEntitlement(name string) Entitlement {
	if perms, ok := entitlements[name]; ok {
		return Entitlement{
			Name:             name,
			AddedPermissions: perms,
		}
	}

	// Panic, because these are hard-coded.
	panic("Entitlement name " + name + " not found!")
}
