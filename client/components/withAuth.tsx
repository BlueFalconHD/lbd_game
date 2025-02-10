import { useContext, useEffect } from "react";
import { useRouter } from "next/router";
import { AuthContext } from "@/context/AuthContext";
import Cookies from "js-cookie";

const withAuth = (WrappedComponent: any, requireAdmin: boolean = false) => {
  const RequiresAuthentication = (props: any) => {
    const { isAuthenticated, isAdmin } = useContext(AuthContext);
    const router = useRouter();

    useEffect(() => {
      const token = Cookies.get("token");
      if (!token) {
        router.push("/login");
      }
    }, []);

    useEffect(() => {
      if (requireAdmin && !isAdmin) {
        router.push("/");
      }
    }, [isAdmin]);

    if (isAuthenticated && (!requireAdmin || isAdmin)) {
      return <WrappedComponent {...props} />;
    } else {
      return null;
    }
  };

  return RequiresAuthentication;
};

export default withAuth;
