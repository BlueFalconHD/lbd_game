import { useContext } from "react";
import { AuthContext } from "@/context/AuthContext";
import Link from "next/link";
import { useRouter } from "next/router";
import Cookies from "js-cookie";

const Navbar = () => {
  const { isAdmin, setAuth } = useContext(AuthContext);
  const router = useRouter();

  const handleLogout = () => {
    Cookies.remove("token");
    setAuth(false, false);
    router.push("/login");
  };

  return (
    <nav className="bg-gray-800 shadow-lg">
      <div className="container mx-auto px-4 py-4 flex justify-between items-center">
        <Link href="/" className="text-xl font-semibold text-gray-100">
          Secret Phrase
        </Link>
        <div className="flex space-x-4">
          {isAdmin && (
            <Link
              href="/admin"
              className="text-gray-100 hover:text-blue-400 transition"
            >
              Admin Dashboard
            </Link>
          )}
          <button
            onClick={handleLogout}
            className="text-gray-100 hover:text-red-400 transition"
          >
            Logout
          </button>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
