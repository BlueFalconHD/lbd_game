import { useState, useEffect } from "react";
import api from "@/utils/api";
import Cookies from "js-cookie";
import { useRouter } from "next/router";
import Link from "next/link";
import withAuth from "@/components/withAuth";
import Layout from "@/components/Layout";

interface Phrase {
  text: string;
  submittedBy: string;
}

interface Verification {
  userId: number;
  username: string;
  confirmedBy: string;
}

interface User {
  username: string;
  isAdmin: boolean;
}

const HomePage = () => {
  const [phrase, setPhrase] = useState<Phrase | null>(null);
  const [verified, setVerified] = useState<boolean>(false);
  const [verifiedBy, setVerifiedBy] = useState<string>("");
  const [users, setUsers] = useState<User[]>([]);
  const [targetUser, setTargetUser] = useState<string>("");
  const [verifications, setVerifications] = useState<Verification[]>([]);
  const [user, setUser] = useState<User | null>(null);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  // Fetch user info
  const fetchUser = async () => {
    try {
      const response = await api.get("/user/status");
      setUser({
        username: response.data.username,
        isAdmin: response.data.is_admin,
      });
    } catch (err: any) {
      if (err.response?.status === 401) {
        router.push("/login");
      }
    }
  };

  // Fetch current phrase
  const fetchPhrase = async () => {
    try {
      const response = await api.get("/user/phrase");
      setPhrase(response.data.phrase);
    } catch (err: any) {
      setPhrase(null);
    }
  };

  // Fetch verification status
  const fetchVerificationStatus = async () => {
    try {
      const response = await api.get("/user/verification_status");
      setVerified(response.data.verified);
      setVerifiedBy(response.data.verified_by);
    } catch (err: any) {
      setVerified(false);
    }
  };

  // Fetch list of users for the dropdown
  const fetchUsers = async () => {
    try {
      const response = await api.get("/user/active_users");
      setUsers(response.data.users);
    } catch (err: any) {
      setUsers([]);
    }
  };

  // Fetch list of verifications
  const fetchVerifications = async () => {
    try {
      const response = await api.get("/user/verifications");
      setVerifications(response.data.verifications);
    } catch (err: any) {
      setVerifications([]);
    }
  };

  useEffect(() => {
    fetchUser();
    fetchPhrase();
    fetchVerificationStatus();
    fetchUsers();
    fetchVerifications();
  }, []);

  const handleVerify = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post("/user/confirm_usage", { username: targetUser });
      setMessage(`Successfully verified ${targetUser}`);
      setError("");
      // Refresh verifications
      fetchVerifications();
    } catch (err: any) {
      setError(err.response?.data?.error || "Verification failed");
      setMessage("");
    }
  };

  const handleLogout = () => {
    Cookies.remove("token");
    router.push("/login");
  };

  return (
    <Layout>
      <div className="max-w-3xl mx-auto">
        <h1 className="text-3xl font-bold mb-6 text-center">Today's Phrase</h1>
        {phrase ? (
          <div className="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
            <p
              className={`text-2xl font-semibold ${
                verified ? "text-green-400" : "text-gray-100"
              } text-center`}
            >
              {phrase.text}
            </p>
            {verified && (
              <p className="mt-2 text-center text-gray-400">
                Verified by <span className="font-medium">{verifiedBy}</span>
              </p>
            )}
          </div>
        ) : (
          <p className="text-center text-gray-400">No phrase submitted yet.</p>
        )}

        <h2 className="text-2xl font-semibold mb-4">Verify a Player</h2>
        {message && <p className="text-green-400 mb-4">{message}</p>}
        {error && <p className="text-red-400 mb-4">{error}</p>}
        <form onSubmit={handleVerify} className="mb-8">
          <div className="flex items-center space-x-4">
            <select
              value={targetUser}
              onChange={(e) => setTargetUser(e.target.value)}
              required
              className="flex-grow bg-gray-700 text-gray-100 p-2 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Select a player</option>
              {(users || [])
                .filter((u) => u.username !== user?.username)
                .map((u) => (
                  <option key={u.username} value={u.username}>
                    {u.username}
                  </option>
                ))}
            </select>
            <button
              type="submit"
              className="bg-blue-600 hover:bg-blue-500 transition text-gray-100 px-4 py-2 rounded"
            >
              Verify
            </button>
          </div>
        </form>

        <h2 className="text-2xl font-semibold mb-4">Verifications Today</h2>
        {verifications ? (
          <ul className="space-y-2">
            {verifications.map((v) => (
              <li
                key={v.userId}
                className="bg-gray-800 p-4 rounded shadow flex justify-between items-center"
              >
                <span>
                  <span className="font-medium">{v.username}</span> verified by{" "}
                  <span className="font-medium">{v.confirmedBy}</span>
                </span>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-gray-400">No verifications yet.</p>
        )}
      </div>
    </Layout>
  );
};

export default withAuth(HomePage);
