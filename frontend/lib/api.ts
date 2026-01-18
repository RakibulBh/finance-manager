"use client";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

export async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;

  const headers = {
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...options.headers,
  };

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Something went wrong");
  }

  return data.data;
}

export async function login(email: string, password: string) {
  const response = await apiRequest<{ token: string; user: { id: string; email: string; family_id: string } }>("/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });

  if (response.token) {
    localStorage.setItem("token", response.token);
    localStorage.setItem("user", JSON.stringify(response.user));
  }

  return response;
}

export async function register(email: string, password: string, familyName: string) {
  const response = await apiRequest<{ token: string; user: { id: string; email: string; family_id: string } }>("/register", {
    method: "POST",
    body: JSON.stringify({ email, password, family_name: familyName }),
  });

  if (response.token) {
    localStorage.setItem("token", response.token);
    localStorage.setItem("user", JSON.stringify(response.user));
  }

  return response;
}

export function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
}
