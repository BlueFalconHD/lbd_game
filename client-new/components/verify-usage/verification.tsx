import { TodayVerification } from "@/lib/api/types";
import { format } from "date-fns";

export interface VerificationProps
  extends React.HTMLAttributes<HTMLDivElement> {
  verification: TodayVerification;
}

export function Verification({ verification, ...props }: VerificationProps) {
  return (
    <div className="flex items-center justify-between py-2" {...props}>
      <div className="flex items-center gap-2">
        <span>{verification.verifier_name}</span>
        <span className="text-muted-foreground">verified</span>
        <span>{verification.verified_name}</span>
      </div>
      <div className="text-sm text-muted-foreground">
        {format(new Date(verification.created_at), "h:mm a")}
      </div>
    </div>
  );
}
