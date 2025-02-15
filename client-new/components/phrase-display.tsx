import useTimeout from "@/hooks/use-timeout";
import { editPhrase, manualReset, unsubmitPhrase } from "@/lib/api/admin";
import {
  canSubmitPhrase,
  getCurrentPhrase,
  submitPhrase,
} from "@/lib/api/phrase";
import { ApiResponse, PhraseResponse } from "@/lib/api/types";
import { Loader2 } from "lucide-react";
import { useContext, useEffect, useState } from "react";
import { toast } from "sonner";
import { EditPhraseDialog } from "./phrase-display/edit-phrase-dialog";
import { ControlledDialog, Dialog } from "./ui/dialog";
import { CurrentPhrase } from "./phrase-display/current-phrase";
import { PhraseInput } from "./phrase-display/phrase-input";
import useInterval from "@/hooks/use-interval";
import { AuthContext } from "@/contexts/AuthContext";
import { Button } from "./ui/button";
import { unixTime5SecondsFromNow } from "@/lib/utils";

export function PhraseDisplay({ ...props }) {
  const [phrase, setPhrase] = useState<string | null>(null);
  const [submittedBy, setSubmittedBy] = useState<string | null>(null);
  const [nextOpenTime, setNextOpenTime] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [newPhrase, setNewPhrase] = useState<string>("");
  const [canSubmit, setCanSubmit] = useState<boolean>(false);
  const { privilege } = useContext(AuthContext);
  const [editDialogOpen, setEditDialogOpen] = useState(false);

  const fetchPhrase = async () => {
    try {
      const phraseResponse = await getCurrentPhrase();
      const submitStatus = await canSubmitPhrase();

      if (!phraseResponse.error) {
        setPhrase(phraseResponse.phrase);
        setSubmittedBy(phraseResponse.submittedBy);
        setNextOpenTime(phraseResponse.next_open_time || null);
      }
      setCanSubmit(submitStatus);
    } catch (error) {
      toast.error("Failed to fetch phrase");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPhrase();
  }, []);

  useEffect(() => {
    if (!nextOpenTime) return;

    const timeUntilOpen = new Date(nextOpenTime).getTime() - Date.now();
    if (timeUntilOpen <= 0) return;

    const timer = setTimeout(() => {
      toast.info("Submission window is now open!");
    }, timeUntilOpen);

    return () => clearTimeout(timer);
  }, [nextOpenTime]);
  useInterval(() => fetchPhrase(), canSubmit ? 500 : 30000);

  const handleEdit = async () => {
    if (!newPhrase.trim()) return;
    try {
      await editPhrase(newPhrase.trim());
      setEditDialogOpen(false);
      toast.success("Phrase edited successfully!");
      setNewPhrase("");
      await fetchPhrase();
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Failed to edit phrase");
    }
  };

  const handleSubmit = async () => {
    if (!newPhrase.trim()) return;
    setSubmitting(true);
    try {
      await submitPhrase(newPhrase.trim());
      toast.success("Phrase submitted successfully!");
      setNewPhrase("");
      await fetchPhrase();
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Failed to submit phrase");
    } finally {
      setSubmitting(false);
    }
  };

  const handleUnsubmit = async () => {
    try {
      await unsubmitPhrase();
      toast.success("Phrase unsubmitted successfully!");
      await fetchPhrase();
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Failed to unsubmit phrase");
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[200px]">
        <Loader2 className="animate-spin h-8 w-8" />
      </div>
    );
  }

  return (
    <ControlledDialog open={editDialogOpen} onOpenChange={setEditDialogOpen}>
      <div className="space-y-4" {...props}>
        {phrase && !canSubmit ? (
          <CurrentPhrase
            phrase={phrase}
            submittedBy={submittedBy!}
            privilege={privilege}
            onEdit={() => {}}
            onUnsubmit={handleUnsubmit}
          />
        ) : canSubmit ? (
          <PhraseInput
            value={newPhrase}
            onChange={setNewPhrase}
            onSubmit={handleSubmit}
            submitting={submitting}
          />
        ) : (
          <div className="flex flex-col justify-center items-center min-h-[200px] space-y-4">
            <Loader2 className="animate-spin h-8 w-8" />
            <p className="text-sm text-muted-foreground">
              No submission window is open.
            </p>
          </div>
        )}
      </div>

      <EditPhraseDialog
        newPhrase={newPhrase}
        setNewPhrase={setNewPhrase}
        onSave={handleEdit}
      />
    </ControlledDialog>
  );
}
