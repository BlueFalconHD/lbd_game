import { TodayVerification } from "@/lib/api/types";
import { Verification } from "./verification";

interface VerifiedUsersProps {
  verifications: TodayVerification[];
}

export function VerifiedUsers({ verifications }: VerifiedUsersProps) {
  return (
    <div className="space-y-1 min-w-[100%] md:min-w-[100%] mx-auto">
      {verifications.map((verification) => (
        <Verification
          key={verification.verification_id}
          verification={verification}
        />
      ))}
    </div>
  );
}
