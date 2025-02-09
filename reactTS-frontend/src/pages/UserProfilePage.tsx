import React, { useState } from "react";
import { useParams } from "react-router-dom";
import { useUser } from "../context/UserContext"; // Import useUser
import { User } from "../models/models";

const UserProfilePage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { fetchUserDetails } = useUser();
  const [user, setUser] = useState<User | undefined>(undefined); 

  React.useEffect(() => {
    const fetchUser = async () => {
      if (id) {
        const userDetails = await fetchUserDetails(parseInt(id));
        setUser(userDetails || undefined);
      }
    };
    fetchUser();
  }, [id, fetchUserDetails]);

  if (!user) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>USER PROFILE</h1>
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
