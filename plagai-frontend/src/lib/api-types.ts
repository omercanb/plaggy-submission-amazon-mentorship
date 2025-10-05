export interface ApiResponse<T = unknown> {
  data?: T;
  message?: string;
  error?: string;
  status?: string;
}

export interface CsrfResponse extends ApiResponse {
  token: string;
}

export interface HealthResponse extends ApiResponse {
  health: string;
}

export interface LoginResponse extends ApiResponse {
  token: string;
}
