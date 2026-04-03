const API = "http://localhost:8080";

function headers() {
  const h = { "Content-Type": "application/json" };
  const token = localStorage.getItem("token");
  if (token) h["Authorization"] = `Bearer ${token}`;
  return h;
}

async function request(method, path, body) {
  const res = await fetch(`${API}${path}`, {
    method,
    headers: headers(),
    body: body ? JSON.stringify(body) : undefined,
  });
  if (res.status === 401) {
    localStorage.removeItem("token");
    window.location.href = "/login";
    return;
  }
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || "falha na requisição");
  return data;
}

export const api = {
  login: (username, password) => request("POST", "/login", { username, password }),
  me: () => request("GET", "/api/me"),

  listUsers: () => request("GET", "/api/users"),
  createUser: (username, password) => request("POST", "/api/users", { username, password }),
  assignRoles: (userId, roleIds) => request("PUT", `/api/users/${userId}/roles`, { role_ids: roleIds }),

  listRoles: () => request("GET", "/api/roles"),
  createRole: (name) => request("POST", "/api/roles", { name }),
  assignPermissions: (roleId, permIds) => request("PUT", `/api/roles/${roleId}/permissions`, { permission_ids: permIds }),

  listPermissions: () => request("GET", "/api/permissions"),
  createPermission: (name) => request("POST", "/api/permissions", { name }),

  listPolicies: () => request("GET", "/api/policies"),
  createPolicy: (name, type, value) => request("POST", "/api/policies", { name, type, value }),
  assignPolicies: (roleId, policyIds) => request("PUT", `/api/roles/${roleId}/policies`, { policy_ids: policyIds }),
};
