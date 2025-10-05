import { User } from "@/types/types";
import { ApiResponse } from "./api-types";

export const API_BASE_URL =
  process.env.ENV === "PROD"
    ? process.env.PUBLIC_API_URL
    : "http://localhost:8080/api/v1";

export async function fetchFromAPI<T = unknown>(
  endpoint: string,
  method: "GET" | "POST" | "PUT" | "DELETE" = "GET",
  body?: unknown,
  user?: User | null
): Promise<ApiResponse<T>> {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    Authorization: user ? `Bearer ${user.token}` : "",
  };

  const config: RequestInit = {
    method,
    headers,
    credentials: "include",
  };

  if (body) {
    config.body = JSON.stringify(body);
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, config);

  if (!response.ok) {
    const errorData: ApiResponse = await response.json().catch(() => ({}));
    throw new Error(errorData.message || "Request failed");
  }

  return response.json() as Promise<ApiResponse<T>>;
}
