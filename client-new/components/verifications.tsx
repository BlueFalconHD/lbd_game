import { useState, useEffect } from "react";
import { TodayVerification } from "@/lib/api/types";
import { getTodayVerifications } from "@/lib/api/verification";
import { PickUser } from "@/components/verify-usage/pick-user";
import { VerifiedUsers } from "@/components/verify-usage/verified-users";
import useInterval from "@/hooks/use-interval";

export function Verifications() {
  const [verifications, setVerifications] = useState<TodayVerification[]>([]);

  const fetchVerifications = async () => {
    const response = await getTodayVerifications();
    if (response.data) {
      setVerifications(response.data);
    }
  };

  useEffect(() => {
    fetchVerifications();
  }, []);

  useInterval(fetchVerifications, 10000);

  return (
    <div className="space-y-6 flex flex-col items-center">
      <PickUser />
      <VerifiedUsers verifications={verifications} />
    </div>
  );
}
