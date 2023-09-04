import { sendPostRequest } from ".";
import { BASE_URL } from "../constants";

export type User = {
  id: string;
  name: string;
  email: string;
  is_guest: boolean;
  is_verified: boolean;
  created_at: string;
  updated_at: string;
};

export type CreateUserParams = {
  name: string;
  email?: string;
  password?: string;
};

export type CreateUserResponse = {
  user: User;
  jwt_token: string;
};

export async function createUser(
  params: CreateUserParams
): Promise<CreateUserResponse> {
  const url = `${BASE_URL}/users`;
  return sendPostRequest<CreateUserResponse>(url, params);
}