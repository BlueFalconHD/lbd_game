import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { KeyboardEvent } from "react";

export function EditPhraseDialog({
  newPhrase,
  setNewPhrase,
  onSave,
}: {
  newPhrase: string;
  setNewPhrase: (phrase: string) => void;
  onSave: () => Promise<void>;
}) {
  const handleKeyDown = async (e: KeyboardEvent) => {
    if (e.key === "Enter") {
      await onSave();
    }
  };

  return (
    <DialogContent className="sm:max-w-[425px]">
      <DialogHeader>
        <DialogTitle>Edit phrase</DialogTitle>
      </DialogHeader>
      <div className="grid gap-4 py-4">
        <div className="grid grid-cols-4 items-center gap-4">
          <Label htmlFor="phrase" className="text-right">
            Phrase
          </Label>
          <Input
            id="phrase"
            value={newPhrase}
            className="col-span-3"
            onChange={(e) => setNewPhrase(e.target.value)}
            onKeyDown={handleKeyDown}
          />
        </div>
      </div>
      <DialogFooter>
        <Button type="submit" onClick={onSave}>
          Save
        </Button>
      </DialogFooter>
    </DialogContent>
  );
}

export function EditPhraseDialogTrigger({
  children,
}: {
  children: React.ReactNode;
}) {
  return <DialogTrigger>{children}</DialogTrigger>;
}
