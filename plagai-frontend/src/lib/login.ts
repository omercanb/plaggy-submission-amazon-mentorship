import { API_BASE_URL } from "./api";
import { LoginResponse } from "./api-types";

export type LoginRequestParams = {
  email: string;
  password: string;
};

export async function loginRequest(
  params: LoginRequestParams
): Promise<LoginResponse> {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };
  const response = await fetch(API_BASE_URL + "/login", {
    method: "POST",
    headers: headers,
    body: JSON.stringify(params),
    credentials: "include",
  });

  if (!response.ok) {
    // Get error message from server
    const errorMessage = await response.text();
    if (errorMessage) {
      throw new Error(errorMessage);
    } else {
      throw new Error("Authentication Failed");
    }
  }

  const token = response.headers.get("token");
  if (!token) {
    const body = await response.json();
    const fallbackToken = body?.token;
    if (!fallbackToken) {
      throw new Error("No Auth Token Found");
    }
    return new Promise<LoginResponse>((resolve) => {
      resolve({
        token: fallbackToken,
        status: "OK",
      });
    });
  }
  return new Promise<LoginResponse>((resolve) => {
    resolve({
      token: token,
      status: "OK",
    });
  }) as Promise<LoginResponse>;
}
