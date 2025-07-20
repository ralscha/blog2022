export interface User {
  id: string;
  email: string;
  name?: string;
  avatar?: string;
  created: string;
  updated: string;
}

export interface AuthData {
  token: string;
  record: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  passwordConfirm: string;
  name?: string;
}

export interface UpdateProfileRequest {
  email?: string;
  name?: string;
}
