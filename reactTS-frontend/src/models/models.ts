// User Model
export interface User {
  id: number;
  username: string;
  email: string;
  password: string;
  bio?: string;
  profilePicture?: string;
  threadsCreated: Thread[];
  role: string;
  participationScore: number;
}

// Thread Model (need to update the thread so that can attach images inside)
export interface Thread {
  id: number;
  title: string;
  content: string;
  attachedImages?: string[];
  postedBy: number;
  categoryId: number; // Single category ID
  createdAt: Date;
  likes: number;
  lastActive: string;
}

// Comment Model (need to update the thread so that can attach images inside)
export interface Comment {
  id: number;
  content: string;
  attachedImages: string[];
  //   posted by which user
  postedBy: number;
  threadId: number;
  createdAt: Date;
  upvotes: number;
  downvotes: number;
}

// Category Model
export interface Category {
  id: number;
  name: string;
  description?: string;
}

// Notification Model
export interface Notification {
  id: number;
  userId: number;
  message: string;
  read: boolean;
  createdAt: Date;
}
