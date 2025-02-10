import { useState } from "react";
import { useRouter } from "next/router";
import api from "@/utils/api";
import Cookies from "js-cookie";
import Link from "next/link";

const LoginPage = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();
  const [error, setError] = useState("");

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await api.post("/login", { username, password });
      console.log(response.data);
      Cookies.set("token", response.data.token);
      router.push("/");
      // Force a reload of page to update the context
      router.reload();
    } catch (err: any) {
      console.log(err.response);
      setError(err.response?.data?.error || "Login failed");
    }
  };

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-3xl font-bold mb-6 text-center">Login</h1>
      <form onSubmit={handleLogin}>
        {error && <p style={{ color: "red" }}>{error}</p>}
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
        <button
          type="submit"
          className="w-full bg-blue-600 hover:bg-blue-500 transition text-gray-100 py-2 rounded mt-4"
        >
          Login
        </button>
      </form>
      <p className="mt-4 text-center text-gray-400">
        Don't have an account?{" "}
        <Link href="/register" className="text-blue-400 hover:underline">
          Register here.
        </Link>
      </p>
    </div>
  );
};

export default LoginPage;
