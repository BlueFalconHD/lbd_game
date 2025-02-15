import { useContext, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { AuthContext } from "@/contexts/AuthContext";
import Cookies from "js-cookie";
import { GetToken } from "@/lib/auth";

const withAuth = (WrappedComponent: any, requirePrivilege: number = 0) => {
  const RequiresAuthentication = (props: any) => {
    const { isAuthenticated, privilege } = useContext(AuthContext);
    const router = useRouter();
    const [loading, setLoading] = useState(true);

    useEffect(() => {
      const token = GetToken();
      if (!token) {
        console.log("No token saved");
        router.push("/login");
      }
      setLoading(false);
    }, []);

    useEffect(() => {
      if (!loading && requirePrivilege > privilege) {
        console.log("Insufficient privilege");
        router.push("/");
      }
    }, [privilege, loading]);

    if (loading) return null; // Add this loading check

    if (isAuthenticated && privilege >= requirePrivilege) {
      return <WrappedComponent {...props} />;
    } else {
      if (!isAuthenticated) {
        console.log("Not authenticated");
      }
      if (privilege < requirePrivilege) {
        console.log("Privilege too low");
      }
      return null;
    }
  };

  return RequiresAuthentication;
};

export default withAuth;
