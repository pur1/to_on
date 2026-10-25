export type User = {
  id: string;
  email: string;
  username: string;
  role: "reader" | "creator" | "admin";
  email_verified: boolean;
  created_at: string;
};

type Session = { user: User; access_token: string };

export class ApiError extends Error {
  status: number;
  code: string;

  constructor(status: number, code: string, message: string) {
    super(message);
    this.status = status;
    this.code = code;
  }
}

// Access token lives in memory only; the refresh token is an httpOnly cookie.
let accessToken: string | null = null;

// Single-flight: concurrent refreshes would reuse a rotated token and revoke the session.
let refreshing: Promise<User> | null = null;

function refreshSession() {
  refreshing ??= request<Session>("/auth/refresh", { method: "POST" }, false)
    .then(storeSession)
    .finally(() => {
      refreshing = null;
    });
  return refreshing;
}

async function request<T>(
  path: string,
  init: RequestInit = {},
  retryOn401 = true,
): Promise<T> {
  const headers = new Headers(init.headers);
  if (init.body) headers.set("Content-Type", "application/json");
  if (accessToken) headers.set("Authorization", `Bearer ${accessToken}`);

  const res = await fetch(`/api/v1${path}`, {
    ...init,
    headers,
    credentials: "include",
  });
  // Expired access token: refresh once and replay the request.
  if (res.status === 401 && retryOn401 && accessToken) {
    accessToken = null;
    await refreshSession();
    return request<T>(path, init, false);
  }
  if (res.status === 204) return undefined as T;

  const body = await res.json().catch(() => null);
  if (!res.ok) {
    throw new ApiError(
      res.status,
      body?.error?.code ?? "unknown",
      body?.error?.message ?? res.statusText,
    );
  }
  return body as T;
}

function storeSession(session: Session) {
  accessToken = session.access_token;
  return session.user;
}

export const authApi = {
  register: (input: { email: string; username: string; password: string }) =>
    request<Session>("/auth/register", {
      method: "POST",
      body: JSON.stringify(input),
    }).then(storeSession),

  login: (input: { identifier: string; password: string }) =>
    request<Session>("/auth/login", {
      method: "POST",
      body: JSON.stringify(input),
    }).then(storeSession),

  refresh: refreshSession,

  logout: () =>
    request<void>("/auth/logout", { method: "POST" }).finally(() => {
      accessToken = null;
    }),

  me: () => request<{ user: User }>("/auth/me").then((r) => r.user),
};
