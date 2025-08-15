export interface AuthInfo {
  motoristaId?: string;
  role?: 'user' | 'admin';
}

const STORAGE_KEY = 'ts_auth';

export function saveAuth(info: AuthInfo) {
  const current = loadAuth();
  const merged = { ...current, ...info };
  localStorage.setItem(STORAGE_KEY, JSON.stringify(merged));
}

export function loadAuth(): AuthInfo {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : {};
  } catch {
    return {};
  }
}

export function isAdmin(): boolean {
  return loadAuth().role === 'admin';
}
