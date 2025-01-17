export interface SecretOutput {
  message: string;
}

export interface SigninInput {
  token: string;
}

export interface CreateTokenInput {
  username: string;
}

export interface CreateTokenOutput {
  token: string;
}
