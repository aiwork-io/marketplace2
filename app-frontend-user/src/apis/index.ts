import axios from "axios";
import { ethers } from "ethers";

import { Auth } from "types/auth";
import { Task, TaskStatus } from "types/task";
import { apis } from "utils/config";

export const login = async (data: { email: string; password: string }) => {
  const res = await axios.post<Auth>(apis.login.url, data);
  return res.data;
};

export const register = async (data: {
  name?: string;
  email: string;
  password: string;
  wallet: string;
}) => {
  if (await signContract(data.wallet, data.email)) {
    const res = await axios.post<Auth>(apis.register.url, data);
    return res.data;
  }
  throw new Error("Verify your wallet failed!");
};

export const getTasks = async (filter?: {
  pageSize?: number;
  pageNumber?: number;
  id?: string;
  createStart?: string;
  createEnd?: string;
  status?: TaskStatus;
}) => {
  const query = {
    _page: filter?.pageNumber || 1,
    _limit: filter?.pageSize || 10,
    id: filter?.id,
    created_at_start: filter?.createStart,
    created_at_end: filter?.createEnd,
    status: filter?.status,
  };
  const res = await axios.get<{ data: Task[]; count: number }>(
    apis.tasks(query).url
  );
  return res.data;
};

export const createTask = async (data: { name: string; category: string }) => {
  const res = await axios.post<Task>(apis.tasks().url, data);
  return res.data;
};

export const uploadImage = async (data: {
  taskId: string;
  filename: string;
}) => {
  const res = await axios.get(apis.uploadImage(data).url);
  return res.data;
};

export const uploadImageToStore = async ({
  file,
  url,
}: {
  file: File;
  url: string;
}) => {
  const res = await fetch(url, {
    method: "PUT",
    headers: { "Content-Type": file.type },
    body: file,
  });
  return res;
};

export const getTaskDetail = async (taskId: string) => {
  const res = await axios.get<Task>(apis.task(taskId).url);
  return res.data;
};

export const updateTaskPayment = async ({
  taskId,
  paymentTxn,
  paymentProof,
}: {
  taskId: string;
  paymentTxn?: string;
  paymentProof?: string;
}) => {
  const res = await axios.patch(apis.updatePayment(taskId).url, {
    ...(paymentTxn ? { payment_txn: paymentTxn } : {}),
    ...(paymentProof ? { payment_proof: paymentProof } : {}),
  });
  return res.data;
};

export const updateProfile = async (data: { name: string; wallet: string }) => {
  if (await signContract(data.wallet, data.name)) {
    const res = await axios.put(apis.profile.url, data);
    return res.data;
  }
  throw new Error("Verify your wallet failed!");
};

export const getProfile = async () => {
  const res = await axios.get<Auth>(apis.getProfile.url);
  return res.data;
};

export const sendResetPasswordEmail = async (email: string) => {
  const res = await axios.post(apis.sendResetPasswordEmail.url, {
    email,
  });
  return res.data;
};

export const verifyAccount = async ({
  state,
  password,
}: {
  state?: string | null;
  password?: string;
}) => {
  const data = {
    state,
    ...(password
      ? {
          payload: {
            password,
          },
        }
      : {}),
  };
  const res = await axios.post<Auth>(apis.verification.url, data);
  return res.data;
};

export const downloadResult = async (taskId: string) => {
  const res = await axios.get<{ url: string }>(apis.downloadResult(taskId).url);
  return res.data;
};

export const signContract = async (wallet: string, content: string) => {
  // @ts-ignore
  const { ethereum } = window;

  if (!ethereum) {
    throw new Error("Make sure you have Metamask installed!");
  }

  const provider = new ethers.providers.Web3Provider(ethereum, "any");
  await provider.send("eth_requestAccounts", []);
  const account = await provider.getSigner();
  const accAddress = await account.getAddress();

  if (accAddress.toLowerCase() !== wallet.toLowerCase()) {
    throw new Error("Please connect your wallet!");
  }

  const message = `Hello ${content}, is this your wallet?`;
  const signature = await account.signMessage(message);
  const actualAddress = ethers.utils.verifyMessage(message, signature);

  return actualAddress?.toLowerCase() === wallet?.toLowerCase();
};

export const getClientTasks = async (filter?: {
  pageSize?: number;
  pageNumber?: number;
  id?: string;
  createStart?: string;
  createEnd?: string;
  status?: TaskStatus;
}) => {
  const query = {
    _page: filter?.pageNumber || 1,
    _limit: filter?.pageSize || 10,
    id: filter?.id,
    created_at_start: filter?.createStart,
    created_at_end: filter?.createEnd,
    status: filter?.status,
  };
  const res = await axios.get<{ data: Task[]; count: number }>(
    apis.clientTasks(query).url
  );
  return res.data;
};
