package helper

type PermissionSet struct {
	byRole map[string]map[string]struct{}
}

func NewPermissionSet(raw map[string][]string) *PermissionSet {
	p := &PermissionSet{byRole: map[string]map[string]struct{}{}}
	for role, perms := range raw {
		p.byRole[role] = map[string]struct{}{}
		for _, perm := range perms {
			p.byRole[role][perm] = struct{}{}
		}
	}
	return p
}
func (p *PermissionSet) Can(role, permission string) bool {
	perms, ok := p.byRole[role]
	if !ok {
		return false
	}
	_, ok = perms[permission]
	return ok
}
