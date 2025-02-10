import { useState, useEffect } from "react";
import api from "@/utils/api";
import { useRouter } from "next/router";
import withAuth from "@/components/withAuth";
import Layout from "@/components/Layout";

// post admin/eliminate_user {username: String}
// post admin/resurrect_user {username: String}
// post admin/approve_user {username: String}
// post admin/unapprove_user {username: String}
// post admin/set_admin {username: String, admin: Boolean}
// post admin/edit_phrase {text: String}
// post admin/reset_game
// get admin/detailed_users -> {
// 	"id":          Number,
// 	"username":    String,
// 	"is_approved": Boolean,
// 	"is_eliminated": Boolean,
// 	"is_admin":    Boolean,
// }

const AdminDashboard = () => {
  const [pendingUsers, setPendingUsers] = useState<{ username: string }[]>([]);
  const [allUsers, setAllUsers] = useState<
    {
      id: number;
      username: string;
      is_approved: boolean;
      is_eliminated: boolean;
      is_admin: boolean;
    }[]
  >([]);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const fetchAllData = async () => {
    try {
      const [pendingResponse, detailedResponse] = await Promise.all([
        api.get("/admin/pending_users"),
        api.get("/admin/detailed_users"),
      ]);
      setPendingUsers(pendingResponse.data.users);
      setAllUsers(detailedResponse.data);
    } catch (err: any) {
      setError("Failed to fetch users");
      if (err.response?.status === 401) {
        router.push("/login");
      }
    }
  };

  useEffect(() => {
    fetchAllData();
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

  const resetGame = async () => {
    try {
      await api.post("/admin/reset_game");
      setMessage("Game reset");
      setError("");
    } catch (err: any) {
      setError("Failed to reset game");
      setMessage("");
    }
  };

  const editPhrase = async (text: string) => {
    try {
      await api.post("/admin/edit_phrase", { text });
      setMessage("Phrase updated");
      setError("");
    } catch (err: any) {
      setError("Failed to update phrase");
      setMessage("");
    }
  };

  const eliminateUser = async (username: string) => {
    try {
      await api.post("/admin/eliminate_user", { username });
      setMessage(`User ${username} eliminated`);
      setError("");
    } catch (err: any) {
      setError(`Failed to eliminate ${username}`);
      setMessage("");
    }
  };

  const resurrectUser = async (username: string) => {
    try {
      await api.post("/admin/resurrect_user", { username });
      setMessage(`User ${username} resurrected`);
      setError("");
    } catch (err: any) {
      setError(`Failed to resurrect ${username}`);
      setMessage("");
    }
  };

  const setAdmin = async (username: string, admin: boolean) => {
    try {
      await api.post("/admin/set_admin", { username, admin });
      setMessage(
        `User ${username} is now ${admin ? "an admin" : "not an admin"}`,
      );
      setError("");
    } catch (err: any) {
      setError(`Failed to set admin status for ${username}`);
      setMessage("");
    }
  };

  const unapproveUser = async (username: string) => {
    try {
      await api.post("/admin/unapprove_user", { username });
      setMessage(`User ${username} unapproved`);
      setError("");
      fetchAllData();
    } catch (err: any) {
      setError(`Failed to unapprove ${username}`);
      setMessage("");
    }
  };

  return (
    <Layout>
      <div className="max-w-5xl mx-auto">
        <h1 className="text-3xl font-bold mb-6 text-center">Admin Dashboard</h1>
        {message && <p className="text-green-400 mb-4">{message}</p>}
        {error && <p className="text-red-400 mb-4">{error}</p>}

        <h2 className="text-2xl font-semibold mb-4">All Users</h2>
        <div className="overflow-x-auto">
          <table className="w-full mb-8">
            <thead>
              <tr className="bg-gray-800">
                <th className="p-3 text-left">Username</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Actions</th>
              </tr>
            </thead>
            <tbody>
              {allUsers.map((user) => (
                <tr key={user.id} className="border-b border-gray-700">
                  <td className="p-3">
                    <span className="font-medium">{user.username}</span>
                  </td>
                  <td className="p-3">
                    <div className="space-y-1">
                      <div
                        className={
                          user.is_approved ? "text-green-400" : "text-red-400"
                        }
                      >
                        {user.is_approved ? "Approved" : "Not Approved"}
                      </div>
                      <div
                        className={
                          user.is_eliminated ? "text-red-400" : "text-green-400"
                        }
                      >
                        {user.is_eliminated ? "Eliminated" : "Active"}
                      </div>
                      <div
                        className={
                          user.is_admin ? "text-blue-400" : "text-gray-400"
                        }
                      >
                        {user.is_admin ? "Admin" : "User"}
                      </div>
                    </div>
                  </td>
                  <td className="p-3">
                    <div className="flex flex-wrap gap-2 justify-center">
                      {user.is_approved ? (
                        <button
                          onClick={() => unapproveUser(user.username)}
                          className="bg-red-600 hover:bg-red-500 px-2 py-1 rounded text-sm"
                        >
                          Unapprove
                        </button>
                      ) : (
                        <button
                          onClick={() => approveUser(user.username)}
                          className="bg-green-600 hover:bg-green-500 px-2 py-1 rounded text-sm"
                        >
                          Approve
                        </button>
                      )}
                      {user.is_eliminated ? (
                        <button
                          onClick={() => resurrectUser(user.username)}
                          className="bg-green-600 hover:bg-green-500 px-2 py-1 rounded text-sm"
                        >
                          Resurrect
                        </button>
                      ) : (
                        <button
                          onClick={() => eliminateUser(user.username)}
                          className="bg-red-600 hover:bg-red-500 px-2 py-1 rounded text-sm"
                        >
                          Eliminate
                        </button>
                      )}
                      <button
                        onClick={() => setAdmin(user.username, !user.is_admin)}
                        className="bg-blue-600 hover:bg-blue-500 px-2 py-1 rounded text-sm"
                      >
                        {user.is_admin ? "Remove Admin" : "Make Admin"}
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <h2 className="text-2xl font-semibold mb-4">Pending Users</h2>
        {pendingUsers.length > 0 ? (
          <ul className="space-y-2">
            {pendingUsers.map((u) => (
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
