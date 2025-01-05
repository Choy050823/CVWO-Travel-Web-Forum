import React from "react";
import { User } from "../models/models";

const UserProfilePage: React.FC<{ user: User }> = ({ user }) => {
  return (
    <div>
      <img src={user.profilePicture !== null ? user.profilePicture! : ""} />
      <ul>
        <li>Username: {user.username}</li>
        <li>User Email: {user.email}</li>
        <li>Bio: {user.bio !== null ? user.bio! : ""}</li>
        <li>Role: {user.role}</li>
      </ul>
    </div>
  );
};

export default UserProfilePage;
