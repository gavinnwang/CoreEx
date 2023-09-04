import { sendPostRequest } from './index';
import { BASE_URL } from '../constants';

export type SignInParams = {
  email: string;
  password: string;
};

type SignInResponse = {
  token: string;
};

export async function signin(params: SignInParams): Promise<SignInResponse> {
  const url = `${BASE_URL}/auth/login`;
  return sendPostRequest<SignInResponse>(url, params);
}