import axios from "axios";

import { Auth } from "types/auth";
import { apis } from "utils/config";

export const login = async (data: { email: string; password: string }) => {
  const res = await axios.post<Auth>(apis.login.url, data);
  return res.data;
};

export const createUser = async (data: {
  name?: string;
  email: string;
  password: string;
  wallet: string;
}) => {
  const res = await axios.post<Auth>(apis.createUser.url, data);
  return res.data;
};
