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
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/signup" element={<SignUpPage />} />
        <Route path="/thread/:id" element={<ThreadDetailsPage />} />
        <Route
          path="/user/:id"
          element={<UserProfilePage user={loginUser} />}
        />
        <Route path="/create" element={<CreateEditThreadPage />} />
        <Route path="/edit/:id" element={<CreateEditThreadPage />} />
        <Route path="/search" element={<SearchPage />} />
        <Route
          path="/notifications"
          element={<NotificationsPage user={loginUser} />}
        />

        <Route path="/admin" element={<AdminDashboardPage />} />
        <Route path="*" element={<ErrorPage />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
