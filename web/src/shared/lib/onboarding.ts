export type UserMode = 'beginner' | 'advanced'

const USER_MODE_KEY = 'crux_user_mode'
export function setUserMode(mode: UserMode) {
  localStorage.setItem(USER_MODE_KEY, mode)
}

export function getPostAuthPath(mode: UserMode | null | undefined): string {
  return mode === 'beginner' ? '/welcome' : '/traders'
}
