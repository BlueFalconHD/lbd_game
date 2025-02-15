import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import { KeyboardEvent } from "react";

function PhraseInput({
  value,
  onChange,
  onSubmit,
  submitting,
}: {
  value: string;
  onChange: (value: string) => void;
  onSubmit: () => Promise<void>;
  submitting: boolean;
}) {
  const handleKeyDown = async (e: KeyboardEvent) => {
    if (e.key === "Enter" && !submitting) {
      await onSubmit();
    }
  };

  // align vertically center
  return (
    <div className="flex gap-2 items-center">
      <Input
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Enter a new phrase..."
        disabled={submitting}
      />
      <Button onClick={onSubmit} disabled={submitting} size="sm">
        {submitting && <Loader2 className="animate-spin" />}
        Submit
      </Button>
    </div>
  );
}

export { PhraseInput };
