import qs from "query-string";

import { TaskStatus } from "types/task";

const filterNonNull = (obj: any) =>
  obj && Object.fromEntries(Object.entries(obj).filter(([_, v]) => v));

export const apis = {
  login: {
    key: "login",
    url: "/auth/login",
  },
  register: {
    key: "register",
    url: "/auth",
  },
  tasks: (query?: {
    _page?: number;
    _limit?: number;
    id?: string;
    created_at_start?: string;
    created_at_end?: string;
    status?: TaskStatus;
  }) => ({
    key: "tasks",
    url: query ? `/tasks?${qs.stringify(filterNonNull(query))}` : "/tasks",
  }),
  task: (taskId?: string) => ({
    key: "task",
    url: `/tasks/${taskId}`,
  }),
  uploadImage: ({
    taskId,
    filename,
  }: {
    taskId: string;
    filename: string;
  }) => ({
    key: "upload-image",
    url: `/tasks/${taskId}/images/upload-url?filename=${filename}`,
  }),
  updatePayment: (taskId: string) => ({
    key: "update-payment",
    url: `/tasks/${taskId}/payment`,
  }),
  profile: {
    key: "profile",
    url: "/users/profile",
  },
  getProfile: {
    key: "get-profile",
    url: "/users/me",
  },
  sendResetPasswordEmail: {
    key: "send-reset-password-email",
    url: "/auth/recovery/password",
  },
  verification: {
    key: "verification",
    url: "/auth/recovery/verification",
  },
  downloadResult: (taskId: string) => ({
    key: "download-result",
    url: `/tasks/${taskId}/results`,
  }),
  clientTasks: (query?: {
    _page?: number;
    _limit?: number;
    id?: string;
    created_at_start?: string;
    created_at_end?: string;
    status?: TaskStatus;
  }) => ({
    key: "tasks",
    url: query
      ? `/client/tasks?${qs.stringify(filterNonNull(query))}`
      : "/client/tasks",
  }),
};

export const procureStatus = [
  {
    value: TaskStatus.PENDING,
    label: "Pending",
  },
  {
    value: TaskStatus.PROCESSING,
    label: "Processing",
  },
  {
    value: TaskStatus.COMPLETED,
    label: "Completed",
  },
];

export const categoryOption = [
  {
    value: "Image Inference",
    label: "Image Inference",
  },
];
