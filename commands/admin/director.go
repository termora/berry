package admin

import "github.com/starshine-sys/bcr"

type directors struct {
	admin *Admin
}

func (d *directors) Check(ctx *bcr.Context) (bool, error) {
	// if the admin check passes, return true
	if t, _ := d.admin.Check(ctx); t {
		return true, nil
	}

	// if there's no member return false
	if ctx.Member == nil {
		return false, nil
	}

	for _, r := range ctx.Member.RoleIDs {
		for _, s := range d.admin.Config.Bot.Support.StaffRoles {
			if r == s {
				return true, nil
			}
		}
	}

	return false, nil
}

func (d *directors) String() string {
	return "Term Director (or Bot Admin)"
}
