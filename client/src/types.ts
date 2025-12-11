export interface User {
  id: number;
  name: string;
  email: string;
}

export interface Subtask {
  id: number;
  task_id: number;
  title: string;
  is_done: boolean;
}

export interface Task {
  id: number;
  user_id: number;
  title: string;
  description: string;
  status: string;
  priority: 'low' | 'medium' | 'high';
  labels: string[];
  reminder_time?: string;
  recurrence_pattern?: string;
  subtasks: Subtask[];
  created_at: string;
}

export interface LoginResponse {
  access_token: string;
  user: { id: number; name: string; email: string };
}
