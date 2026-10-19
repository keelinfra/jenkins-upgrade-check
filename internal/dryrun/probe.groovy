// Written by jenkins-upgrade-check into the rehearsal container only. It
// records what this Jenkins loaded and what it could not, in a JSON file the
// tool reads back with "docker cp". It changes nothing.
import groovy.json.JsonOutput
import hudson.init.InitMilestone
import jenkins.model.Jenkins

// where names a saveable thing the way an administrator would look for it.
def where = { k ->
    if (k instanceof hudson.model.Item) return k.fullName
    if (k instanceof hudson.model.Run) return k.parent.fullName + ' #' + k.number
    if (k instanceof jenkins.model.Jenkins) return '(global configuration)'
    if (k instanceof hudson.model.Node) return '(node ' + k.nodeName + ')'
    return String.valueOf(k)
}

def collect = { boolean done ->
    def out = [completed: done]
    def errors = []
    try {
        def j = Jenkins.get()
        out.version = Jenkins.VERSION
        out.java = System.getProperty('java.version')
        def pm = j.pluginManager
        out.plugins = pm.plugins.collect { p ->
            [name: p.shortName, version: p.version, active: p.isActive(), enabled: p.isEnabled()]
        }
        out.failed = pm.failedPlugins.collect { f ->
            def chain = []
            def c = f.cause
            while (c != null && chain.size() < 8) {
                chain << (c.getClass().name + ': ' + (c.message ?: ''))
                if (c.cause == null || c.cause.is(c)) break
                c = c.cause
            }
            [name: f.name, causes: chain]
        }
        def root = j.rootDir.toPath()
        try {
            out.items = j.getAllItems(hudson.model.Item).collect { root.relativize(it.rootDir.toPath()).toString() }
        } catch (Throwable t) { errors << ('items: ' + t) }
        try {
            def odm = j.getExtensionList(hudson.diagnosis.OldDataMonitor).get(0)
            out.unreadable = odm.data.findAll { k, v -> v?.extra }.collect { k, v ->
                [item: where(k), path: (k instanceof hudson.model.Item) ? root.relativize(k.rootDir.toPath()).toString() : '', problem: v.extra]
            }
        } catch (Throwable t) { errors << ('old data: ' + t) }
        if (done) {
            try {
                out.monitors = j.administrativeMonitors.findAll { m ->
                    try { m.isEnabled() && m.isActivated() } catch (Throwable t) { false }
                }.collect { m -> [id: m.id, name: m.displayName] }
            } catch (Throwable t) { errors << ('monitors: ' + t) }
        }
    } catch (Throwable t) {
        errors << t.toString()
    }
    out.probeErrors = errors
    def f = new File(Jenkins.get().rootDir, 'jenkins-upgrade-check-probe.json')
    def tmp = new File(f.path + '.tmp')
    tmp.text = JsonOutput.toJson(out)
    tmp.renameTo(f)
}

// asSystem runs body with full rights. Without them, getAllItems leaves out
// every item the anonymous user cannot read, which on a secured controller
// is all of them.
def asSystem = { Closure body ->
    def ctx
    try {
        ctx = hudson.security.ACL.as2(hudson.security.ACL.SYSTEM2)
    } catch (MissingPropertyException | MissingMethodException e) {
        ctx = hudson.security.ACL.as(hudson.security.ACL.SYSTEM) // before 2.266
    }
    try {
        body()
    } finally {
        ctx.close()
    }
}

// Write once now, in case Jenkins never finishes starting, and again once
// it has.
asSystem { collect(false) }
Thread.start('jenkins-upgrade-check-probe') {
    def deadline = System.currentTimeMillis() + 30 * 60 * 1000
    while (Jenkins.get().initLevel != InitMilestone.COMPLETED && System.currentTimeMillis() < deadline) {
        sleep(500)
    }
    sleep(3000)
    asSystem { collect(true) }
}
