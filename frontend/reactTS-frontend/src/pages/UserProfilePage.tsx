import React from "react";
import { User } from "../models/models";

interface UserProfilePageProps {
  user: User;
}

const UserProfilePage: React.FC<UserProfilePageProps> = ({ user }) => {
  return (
    <div>
      <img src={user.profilePicture || ""} alt="Profile" />
      <ul>
        <li>Username: {user.username}</li>
        <li>User Email: {user.email}</li>
        <li>Bio: {user.bio || "No bio available"}</li>
        <li>Role: {user.role}</li>
      </ul>
    </div>
  );
};

export default UserProfilePage;
