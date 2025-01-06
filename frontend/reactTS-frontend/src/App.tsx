import { Routes, Route, BrowserRouter } from "react-router-dom";
import HomePage from "./pages/HomePage";
import { User } from "./models/models";
import LoginPage from "./pages/LoginPage";
import SignUpPage from "./pages/SignUpPage";
import ErrorPage from "./pages/ErrorPage";
import UserProfilePage from "./pages/UserProfilePage";
import ThreadDetailsPage from "./pages/ThreadDetailsPage";
import AdminDashboardPage from "./pages/AdminDashboardPage";
import CreateEditThreadPage from "./pages/CreateEditThreadPage";
import SearchPage from "./pages/SearchPage";
import NotificationsPage from "./pages/NotificationsPage";
import ProtectedRoute from "./components/ProtectedRoute";
import AdminRoute from "./components/AdminRoute";

const App = () => {
  const loginUser: User = {
    id: 1,
    username: "john_doe",
    password: "",
    profilePicture: "https://randomuser.me/api/portraits",
    email: "john_doe@gmail.com",
    bio: "Hello, I'm John Doe",
    role: "admin",
    threadsCreated: [],
    participationScore: 0,
  };

  return (
    <BrowserRouter>
      <Routes>
        {/* General Users */}
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/signup" element={<SignUpPage />} />
        <Route path="/search" element={<SearchPage />} />
        <Route path="*" element={<ErrorPage />} />

        {/* Login Users */}
        <Route
          path="/thread/:id"
          element={
            <ProtectedRoute user={loginUser}>
              <ThreadDetailsPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/user/:id"
          element={
            <ProtectedRoute user={loginUser}>
              <UserProfilePage user={loginUser} />
            </ProtectedRoute>
          }
        />
        <Route
          path="/create"
          element={
            <ProtectedRoute user={loginUser}>
              <CreateEditThreadPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/edit/:id"
          element={
            <ProtectedRoute user={loginUser}>
              <CreateEditThreadPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/notifications"
          element={
            <ProtectedRoute user={loginUser}>
              <NotificationsPage user={loginUser} />
            </ProtectedRoute>
          }
        />

        {/* Admin Only */}
        <Route
          path="/admin"
          element={
            <AdminRoute user={loginUser}>
              <AdminDashboardPage />
            </AdminRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
