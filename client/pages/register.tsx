import { useState } from "react";
import { useRouter } from "next/router";
import api from "@/utils/api";
import Link from "next/link";

const RegisterPage = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const router = useRouter();
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }
    try {
      await api.post("/register", { username, password });
      setMessage("Registration successful! Awaiting admin approval.");
      // Redirect to login page after a delay
      setTimeout(() => {
        router.push("/login");
      }, 3000);
    } catch (err: any) {
      console.log(err.response);
      setError(err.response?.data?.error || "Registration failed");
    }
  };

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-3xl font-bold mb-6 text-center">Register</h1>
      <form onSubmit={handleRegister}>
        {error && <p style={{ color: "red" }}>{error}</p>}
        {message && <p style={{ color: "green" }}>{message}</p>}
        <div className="mb-4">
          <label className="block mb-1">Username:</label>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full bg-gray-700 text-gray-100 p-2 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block mb-1">Password:</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full bg-gray-700 text-gray-100 p-2 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block mb-1">Confirm Password:</label>
          <input
            type="password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            className="w-full bg-gray-700 text-gray-100 p-2 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>
        <button
          type="submit"
          className="w-full bg-blue-600 hover:bg-blue-500 transition text-gray-100 py-2 rounded mt-4"
        >
          Register
        </button>
      </form>
      <p className="mt-4 text-center text-gray-400">
        Already have an account?{" "}
        <Link href="/login" className="text-blue-400 hover:underline">
          Sign in here.
        </Link>
      </p>
    </div>
  );
};

export default RegisterPage;
