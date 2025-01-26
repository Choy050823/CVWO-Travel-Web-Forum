import React from "react";
import { useParams } from "react-router-dom";
import { useUser } from "../context/UserContext"; // Import useUser

const UserProfilePage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { userDetails, fetchUserDetails } = useUser();

  React.useEffect(() => {
    if (id) {
      fetchUserDetails(parseInt(id));
    }
  }, [id, fetchUserDetails]);

  if (!userDetails) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>USER PROFILE</h1>
      <img src={userDetails.profilePicture || ""} alt="Profile" />
      <ul>
        <li>Username: {userDetails.username}</li>
        <li>User Email: {userDetails.email}</li>
        <li>Bio: {userDetails.bio || "No bio available"}</li>
        <li>Role: {userDetails.role}</li>
      </ul>
    </div>
  );
};

export default UserProfilePage;
