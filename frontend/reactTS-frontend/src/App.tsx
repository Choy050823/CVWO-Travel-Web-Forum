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
import { useEffect, useState } from "react";
import BASE_URL from "./config";

const App = () => {
  const [user, setUser] = useState<User | null>(null);

  // Fetch user data after login
  useEffect(() => {
    const fetchUserData = async () => {
      const token = localStorage.getItem("token");
      if (token) {
        try {
          const response = await fetch(`${BASE_URL}/api/users/me`, {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });
          if (response.ok) {
            const userData = await response.json();
            setUser(userData);
          }
        } catch (error) {
          console.error("Error fetching user data:", error);
        }
      }
    };

    fetchUserData();
  }, []);

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
            <ProtectedRoute user={user}>
              <ThreadDetailsPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/user/:id"
          element={
            <ProtectedRoute user={user}>
              <UserProfilePage user={user!} />
            </ProtectedRoute>
          }
        />
        <Route
          path="/create"
          element={
            <ProtectedRoute user={user}>
              <CreateEditThreadPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/edit/:id"
          element={
            <ProtectedRoute user={user}>
              <CreateEditThreadPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/notifications"
          element={
            <ProtectedRoute user={user}>
              <NotificationsPage user={user!} />
            </ProtectedRoute>
          }
        />

        {/* Admin Only */}
        <Route
          path="/admin"
          element={
            <AdminRoute user={user}>
              <AdminDashboardPage />
            </AdminRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
