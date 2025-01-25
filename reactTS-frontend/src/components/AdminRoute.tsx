import { Navigate } from "react-router-dom";
import { User } from "../models/models";

const AdminRoute: React.FC<{ user: User | null; children: JSX.Element }> = ({
  user,
  children,
}) => {
  if (!user || user.role !== "admin") {
    return <Navigate to="/" />;
  }
  return children;
};

export default AdminRoute;
