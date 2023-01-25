import { TaskStatus } from "types/task";
import { ACCESS_TOKEN_KEY } from "./constants";

export const isLoggedIn = () => {
  return !!localStorage.getItem(ACCESS_TOKEN_KEY);
};

export const setAuth = (accessToken: string) => {
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
};

export const removeAuth = () => {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
};

export const logout = () => {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  window.location.href = "/login";
};

export const delay = async (t: number) =>
  new Promise((resolve) => setTimeout(resolve, t));

export const generateTaskStatus = (status: TaskStatus) => {
  switch (status) {
    case TaskStatus.PENDING:
      return "Pending";
    case TaskStatus.PROCESSING:
      return "Processing";
    case TaskStatus.COMPLETED:
      return "Completed";
    default:
      return "";
  }
};
