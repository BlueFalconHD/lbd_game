import { DialogTrigger } from "@/components/ui/dialog";
import { SlimButton } from "@/components/ui/slim-button";
import { Pencil, Trash2 } from "lucide-react";

function CurrentPhrase({
  phrase,
  submittedBy,
  nextOpenTime,
  privilege,
  onEdit,
  onUnsubmit,
}: {
  phrase: string;
  submittedBy: string;
  nextOpenTime?: string;
  privilege: number;
  onEdit: () => void;
  onUnsubmit: () => Promise<void>;
}) {
  return (
    <div className="text-center space-y-2 max-w-full w-full min-w-full">
      <div className="max-w-[90%] md:max-w-[600px] mx-auto">
        <h2 className="text-3xl font-bold break-words whitespace-pre-wrap">
          {phrase}
        </h2>
      </div>
      <p className="text-sm text-muted-foreground">
        Submitted by: {submittedBy}
      </p>
      {nextOpenTime && (
        <p className="text-sm text-muted-foreground">
          Next window opens at: {new Date(nextOpenTime).toLocaleString()}
        </p>
      )}
      {privilege > 0 && (
        <div className="flex justify-center items-center gap-2">
          <div className="flex gap-2">
            <DialogTrigger asChild>
              <SlimButton variant="outline">
                <Pencil />
                Edit
              </SlimButton>
            </DialogTrigger>
            <SlimButton variant="outline" onClick={onUnsubmit}>
              <Trash2 />
              Unsubmit
            </SlimButton>
          </div>
        </div>
      )}
    </div>
  );
}

export { CurrentPhrase };
