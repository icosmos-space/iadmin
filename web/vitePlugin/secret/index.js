export function AddSecret(secret, projectName) {
  if (!secret) {
    secret = ''
  }
  if (!projectName) {
    projectName = ''
  }
  global['iadmin-secret'] = secret
  global['iadmin-project-name'] = projectName
}
