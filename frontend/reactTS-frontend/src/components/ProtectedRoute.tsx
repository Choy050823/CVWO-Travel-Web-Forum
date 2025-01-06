import { Navigate } from "react-router-dom";
import { User } from "../models/models";

const ProtectedRoute: React.FC<{
  user: User | null;
  children: JSX.Element;
}> = ({ user, children }) => {
  if (!user) {
    return <Navigate to="/" />;
  }
  return children;
};

export default ProtectedRoute;
