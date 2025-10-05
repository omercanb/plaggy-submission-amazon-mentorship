/**
 *  UNUSED SINCE I NUKED CSRF FROM BACKEND. KEEP THIS IN CASE WE NEED IT ~brtcrt
 * @param API_BASE_URL
 * @returns
 */

export async function getCsrfToken(API_BASE_URL: string): Promise<string> {
  const response = await fetch(API_BASE_URL + "/csrf-token", {
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Failed to get CSRF token");
  }

  const token = response.headers.get("X-Csrf-Token");
  if (!token) {
    throw new Error("CSRF token not found in response");
  }

  return token;
}
