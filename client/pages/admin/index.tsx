import { useState, useEffect } from "react";
import api from "@/utils/api";
import { useRouter } from "next/router";
import withAuth from "@/components/withAuth";
import Layout from "@/components/Layout";

const AdminDashboard = () => {
  const [users, setUsers] = useState<{ username: string }[]>([]);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const fetchPendingUsers = async () => {
    try {
      const response = await api.get("/admin/pending_users");
      setUsers(response.data.users);
    } catch (err: any) {
      setUsers([]);
      setError("Failed to fetch pending users");
      if (err.response?.status === 401) {
        router.push("/login");
      }
    }
  };

  useEffect(() => {
    fetchPendingUsers();
  }, []);

  const approveUser = async (username: string) => {
    try {
      await api.post("/admin/approve_user", { username });
      setMessage(`User ${username} approved`);
      setError("");
      // Refresh the user list
      fetchPendingUsers();
    } catch (err: any) {
      setError(`Failed to approve ${username}`);
      setMessage("");
    }
  };

  return (
    <Layout>
      <div className="max-w-3xl mx-auto">
        <h1 className="text-3xl font-bold mb-6 text-center">Admin Dashboard</h1>
        {message && <p className="text-green-400 mb-4">{message}</p>}
        {error && <p className="text-red-400 mb-4">{error}</p>}
        <h2 className="text-2xl font-semibold mb-4">Pending Users</h2>
        {users ? (
          <ul className="space-y-2">
            {users.map((u) => (
              <li
                key={u.username}
                className="bg-gray-800 p-4 rounded shadow flex justify-between items-center"
              >
                <span className="font-medium">{u.username}</span>
                <button
                  onClick={() => approveUser(u.username)}
                  className="bg-green-600 hover:bg-green-500 transition text-gray-100 px-4 py-2 rounded"
                >
                  Approve
                </button>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-gray-400">No pending users.</p>
        )}
      </div>
    </Layout>
  );
};

export default withAuth(AdminDashboard, true);
