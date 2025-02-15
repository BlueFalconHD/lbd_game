import { useEffect, useState } from "react";
import { LazySelect } from "@/components/ui/lazy-select";
import {
  getTodayVerifications,
  getUnverifiedUsers,
  verifyUser,
} from "@/lib/api/verification";
import { UnverifiedUsersResponse } from "@/lib/api/types";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";

export function PickUser() {
  const [error, setError] = useState<string | null>(null);
  const [selected, setSelected] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const loadOptions = async () => {
    try {
      const response = await getUnverifiedUsers();

      if (response.error || !response.data) {
        setError(response.error || "An unknown error occurred");
        return [];
      }

      return response.data.map((user) => ({
        value: user.id.toString(),
        label: user.username,
      }));
    } catch (err: any) {
      setError(err.message);
      return [];
    }
  };

  function safeSetSelected(value: string | null) {
    if (value === null) {
      setSelected(null);
    } else {
      setSelected(value);
    }
  }

  const onSubmit = async () => {
    if (!selected) {
      setError("Please select a user to verify");
      return;
    }

    setSubmitting(true);
    try {
      await verifyUser(parseInt(selected));
      // setMessage("User verified successfully");
      setSelected(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  const fetchVerifications = async () => {
    const v = await getTodayVerifications();
    console.log(v);
  };

  useEffect(() => {
    fetchVerifications();
  }, []);

  return (
    <div className="flex gap-2 items-center">
      <LazySelect
        value={selected}
        onChange={(value) => safeSetSelected(value)}
        loadOptions={loadOptions}
        noItemsMessage="No unverified users"
        placeholder="Select a user to verify"
        className="w-[280px]"
      />
      <Button onClick={onSubmit} disabled={submitting} size="sm">
        {submitting && <Loader2 className="animate-spin" />}
        Verify
      </Button>
    </div>
  );
}

export default PickUser;
