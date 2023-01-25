export type Task = {
  id: string;
  name: string;
  status: TaskStatus;
  user_id: string;
  created_at: string;
  assets: TaskAsset[];
  completed_at: string;
};

export enum TaskStatus {
  PENDING = 0,
  PROCESSING = 100,
  COMPLETED = 200,
}

type TaskAsset = {
  id: string;
  file_url: string;
  results: string;
};

export type Result = {
  source: string;
  data: {
    object: {
      bbox: number[];
      category: string;
      score: number;
    }[];
  };
};
