package plan

import (
	"github.com/keelinfra/jenkins-upgrade-check/internal/uc"
	"github.com/keelinfra/jenkins-upgrade-check/internal/version"
)

// split is a plugin Jenkins split out of core: core's
// jenkins/split-plugins.txt. A plugin built for a core no newer than the
// split still expects it, and Jenkins gives it the split plugin as an
// implied dependency. Jenkins used to install these from copies inside the
// war; 2.574 and 2.577 removed the copies, so a controller without the
// update center has to bring them itself.
type split struct{ name, when, min string }

// splits are the plugins the war carried until 2.574 and 2.577. Older
// splits (maven-plugin, subversion, ldap and the like) were dropped from
// the war long before and are not implied any more in practice.
var splits = []split{
	{"mailer", "1.493", "1.2"},
	{"matrix-auth", "1.535", "1.0.2"},
	{"antisamy-markup-formatter", "1.553", "1.0"},
	{"matrix-project", "1.561", "1.0"},
	{"junit", "1.577", "1.0"},
	{"bouncycastle-api", "2.16", "2.16.0"},
	{"command-launcher", "2.86", "1.0"},
	{"jdk-tool", "2.112", "1.0"},
	{"jaxb", "2.163", "2.3.0"},
	{"trilead-api", "2.184", "1.0.4"},
	{"sshd", "2.281", "3.236.ved5e1b_cb_50b_2"},
	{"javax-activation-api", "2.330", "1.2.0-2"},
	{"javax-mail-api", "2.330", "1.6.2-5"},
	{"instance-identity", "2.356", "3.1"},
}

// splitCycles are the pairs core's jenkins/split-plugin-cycles.txt exempts:
// the first plugin does not get the second as an implied dependency.
var splitCycles = map[string]bool{
	"credentials matrix-auth":                   true,
	"script-security antisamy-markup-formatter": true,
	"script-security bouncycastle-api":          true,
	"script-security command-launcher":          true,
	"script-security matrix-auth":               true,
	"script-security matrix-project":            true,
	"script-security jaxb":                      true,
	"mailer jaxb":                               true,
	"matrix-auth jaxb":                          true,
	"antisamy-markup-formatter jaxb":            true,
	"matrix-project jaxb":                       true,
	"junit jaxb":                                true,
	"bouncycastle-api jaxb":                     true,
	"command-launcher jaxb":                     true,
	"jdk-tool jaxb":                             true,
	"javax-activation-api javax-mail-api":       true,
	"javax-activation-api sshd":                 true,
	"javax-mail-api sshd":                       true,
	"bouncycastle-api instance-identity":        true,
	"bouncycastle-api sshd":                     true,
	"javax-activation-api instance-identity":    true,
	"javax-mail-api instance-identity":          true,
	"maven-plugin jaxb":                         true,
	"subversion jaxb":                           true,
	"cvs jaxb":                                  true,
	"ant jaxb":                                  true,
	"javadoc jaxb":                              true,
	"external-monitor-job jaxb":                 true,
	"ldap jaxb":                                 true,
	"pam-auth jaxb":                             true,
}

// implied returns the split plugins a release expects on this core: those
// split out by the core it runs on, after the core it was built for.
func implied(r *uc.Release, core string) []uc.Dependency {
	var out []uc.Dependency
	for _, s := range splits {
		if s.name == r.Name || splitCycles[r.Name+" "+s.name] || version.Less(core, s.when) {
			continue
		}
		if r.RequiredCore != "" && version.Less(s.when, r.RequiredCore) {
			continue
		}
		out = append(out, uc.Dependency{Name: s.name, Version: s.min})
	}
	return out
}
