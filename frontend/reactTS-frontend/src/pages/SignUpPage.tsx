import { useNavigate } from "react-router-dom";
import { useState } from "react";
import BASE_URL from "../config";

const SignUpPage: React.FC = () => {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const navigate = useNavigate();

  const handleSignUp = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const response = await fetch(`${BASE_URL}/api/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, email, password }),
      });

      if (!response.ok) {
        throw new Error("Sign up failed");
      }

      navigate("/login"); // Redirect to login page
    } catch (error) {
      console.error("Error during sign up:", error);
    }
  };

  return (
    <div>
      <h2>Sign Up</h2>
      <form onSubmit={handleSignUp}>
        <div>
          <label htmlFor="username">Username:</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="email">Email:</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit">Sign Up</button>
      </form>
    </div>
  );
};

export default SignUpPage;
