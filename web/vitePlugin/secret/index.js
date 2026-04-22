export function AddSecret(secret) {
  if (!secret) {
    secret = ''
  }
  global['iadmin-secret'] = secret
}
