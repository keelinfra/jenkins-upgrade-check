// Gives a fresh controller the kind of configuration a real one collects
// through the UI over the years, then deletes itself. Uses only APIs that
// have been stable since well before 2.346.3.
import com.cloudbees.hudson.plugins.folder.Folder
import com.cloudbees.plugins.credentials.CredentialsScope
import com.cloudbees.plugins.credentials.SystemCredentialsProvider
import com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl
import hudson.model.FreeStyleProject
import hudson.model.ListView
import hudson.model.ParametersDefinitionProperty
import hudson.model.StringParameterDefinition
import hudson.security.GlobalMatrixAuthorizationStrategy
import hudson.security.HudsonPrivateSecurityRealm
import hudson.slaves.DumbSlave
import hudson.slaves.JNLPLauncher
import hudson.slaves.RetentionStrategy
import hudson.tasks.Shell
import hudson.triggers.TimerTrigger
import hudson.util.Secret
import jenkins.model.Jenkins
import org.jenkinsci.plugins.plaincredentials.impl.StringCredentialsImpl
import org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition
import org.jenkinsci.plugins.workflow.job.WorkflowJob

def j = Jenkins.get()
if (j.getItem('platform') != null) {
  return
}

def realm = new HudsonPrivateSecurityRealm(false)
realm.createAccount('admin', System.getenv('E2E_ADMIN_PASSWORD') ?: 'e2e-admin')
j.setSecurityRealm(realm)
def auth = new GlobalMatrixAuthorizationStrategy()
auth.add(Jenkins.ADMINISTER, 'admin')
auth.add(Jenkins.READ, 'authenticated')
j.setAuthorizationStrategy(auth)
j.setSystemMessage('e2e: a controller upgraded by following the runbook')
j.setNumExecutors(2)

def agent = new DumbSlave('agent1', '/home/jenkins/agent', new JNLPLauncher(true))
agent.setLabelString('linux')
agent.setNumExecutors(2)
agent.setRetentionStrategy(new RetentionStrategy.Always())
j.addNode(agent)
// An agent nobody runs any more: offline before the upgrade and after.
def gone = new DumbSlave('agent-retired', '/home/jenkins/agent', new JNLPLauncher(true))
gone.setLabelString('retired')
gone.setRetentionStrategy(new RetentionStrategy.Always())
j.addNode(gone)

def store = SystemCredentialsProvider.getInstance()
store.getCredentials().add(new UsernamePasswordCredentialsImpl(CredentialsScope.GLOBAL, 'deploy-user', 'deploy account', 'deployer', 'not-a-real-password'))
store.getCredentials().add(new StringCredentialsImpl(CredentialsScope.GLOBAL, 'npm-token', 'npm token', Secret.fromString('not-a-real-token')))
store.save()

def platform = j.createProject(Folder, 'platform')
def team = platform.createProject(Folder, 'team-a')

def pipeline = { parent, name, script ->
  def p = parent.createProject(WorkflowJob, name)
  p.setDefinition(new CpsFlowDefinition(script, true))
  p.save()
  p
}
pipeline(platform, 'declarative', '''pipeline {
  agent any
  options { timestamps(); buildDiscarder(logRotator(numToKeepStr: '20')) }
  stages {
    stage('build') { steps { sh 'echo building on $NODE_NAME' } }
    stage('test') { steps { sh 'echo testing' } }
  }
}''')
pipeline(platform, 'on-agent', '''node('linux') {
  sh 'java -version 2>&1 | head -1; echo built on $NODE_NAME'
}''')
pipeline(team, 'deploy', '''node {
  withCredentials([usernamePassword(credentialsId: 'deploy-user', usernameVariable: 'U', passwordVariable: 'P'),
                   string(credentialsId: 'npm-token', variable: 'T')]) {
    sh 'echo deploying as $U'
  }
}''')
def nightly = pipeline(team, 'nightly', '''pipeline {
  agent any
  parameters { string(name: 'TARGET', defaultValue: 'staging') }
  stages { stage('run') { steps { echo "nightly for ${params.TARGET}" } } }
}''')
nightly.addTrigger(new TimerTrigger('H H(0-3) * * *'))
nightly.save()

def legacy = j.createProject(FreeStyleProject, 'legacy-build')
legacy.addProperty(new ParametersDefinitionProperty(new StringParameterDefinition('VERSION', '1.0', 'version to build')))
legacy.getBuildersList().add(new Shell('echo building version $VERSION'))
legacy.save()

def view = new ListView('Platform', j)
view.setRecurse(true)
view.setIncludeRegex('platform/.*')
j.addView(view)

j.save()
new File(j.getRootDir(), 'init.groovy.d/seed.groovy').delete()
