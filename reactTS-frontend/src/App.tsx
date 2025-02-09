import { Routes, Route } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import { UserProvider } from "./context/UserContext"; // Import UserProvider
import { ThreadProvider } from "./context/ThreadContext";
import Navbar from "./components/Navbar";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import SignUpPage from "./pages/SignUpPage";
import ErrorPage from "./pages/ErrorPage";
import UserProfilePage from "./pages/UserProfilePage";
import ThreadDetailsPage from "./pages/ThreadDetailsPage";
import AdminDashboardPage from "./pages/AdminDashboardPage";
import CreateEditThreadPage from "./pages/CreateEditThreadPage";
import SearchPage from "./pages/SearchPage";
import ProtectedRoute from "./components/ProtectedRoute";
import AdminRoute from "./components/AdminRoute";
import { CommentProvider } from "./context/CommentContext";

const App = () => {
  return (
    // <BrowserRouter>
      <AuthProvider>
        <UserProvider>
          {/* Wrap with UserProvider */}
          <ThreadProvider>
            <CommentProvider>
              <Navbar />

              <Routes>
                {/* General Users */}
                <Route path="/" element={<HomePage />} />
                <Route path="/login" element={<LoginPage />} />
                <Route path="/signup" element={<SignUpPage />} />
                <Route path="/search" element={<SearchPage />} />
                <Route path="/threads/:id" element={<ThreadDetailsPage />} />

                <Route path="*" element={<ErrorPage />} />

                {/* Login Users */}
                <Route
                  path="/user/:id"
                  element={
                    <ProtectedRoute>
                      <UserProfilePage />
                    </ProtectedRoute>
                  }
                />

                <Route
                  path="/create-thread"
                  element={
                    <ProtectedRoute>
                      <CreateEditThreadPage />
                    </ProtectedRoute>
                  }
                />
                <Route
                  path="/edit-thread/:id"
                  element={
                    <ProtectedRoute>
                      <CreateEditThreadPage />
                    </ProtectedRoute>
                  }
                />

                {/* Admin Only */}
                <Route
                  path="/admin"
                  element={
                    <AdminRoute>
                      <AdminDashboardPage />
                    </AdminRoute>
                  }
                />
              </Routes>
            </CommentProvider>
          </ThreadProvider>
        </UserProvider>
      </AuthProvider>
    // </BrowserRouter>
  );
};

export default App;
