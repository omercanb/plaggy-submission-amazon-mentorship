"use client";

import { useEffect, useMemo, useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Stack,
  Autocomplete,
  Chip,
  Snackbar,
  CircularProgress,
} from "@mui/material";
import { fetchFromAPI } from "@/lib/api";
import { useAuth } from "@/providers/AuthProvider";

export type Section = { id: string; name: string };

type Props = {
  open: boolean;
  onClose: () => void;
};

export default function AssignHomeworkDialog({ open, onClose }: Props) {
  const { user } = useAuth();

  const [sections, setSections] = useState<Section[]>([]);
  const [loadingSections, setLoadingSections] = useState(false);

  const [selected, setSelected] = useState<Section[]>([]);
  const [title, setTitle] = useState("");
  const [due, setDue] = useState<string>("");
  const [desc, setDesc] = useState("");
  const [busy, setBusy] = useState(false);
  const [toast, setToast] = useState<{
    open: boolean;
    msg: string;
    error?: boolean;
  }>({ open: false, msg: "" });

  useEffect(() => {
    if (!open) return;
    const load = async () => {
      setLoadingSections(true);
      try {
        const res = await fetchFromAPI<Section[]>(
          "/sections",
          "GET",
          null,
          user
        );
        setSections(res.data || []);
      } catch (e) {
        setToast({ open: true, msg: "Failed to load sections", error: true });
      } finally {
        setLoadingSections(false);
      }
    };
    load();
  }, [open, user]);

  const sectionIDs = useMemo(() => selected.map((s) => s.id), [selected]);

  const handleSubmit = async () => {
    if (!title.trim() || !due || sectionIDs.length === 0) {
      setToast({
        open: true,
        msg: "Please choose sections, title, and due date.",
        error: true,
      });
      return;
    }
    setBusy(true);
    try {
      const dueAtISO = new Date(due).toISOString();
      const uintSectionIDs: number[] = [];
      sectionIDs.forEach((id) => {
        uintSectionIDs.push(Number.parseInt(id));
      });
      const payload = {
        sectionIDs: uintSectionIDs,
        title: title.trim(),
        dueAt: dueAtISO,
      };

      const res = await fetchFromAPI(
        "/create_homeworks",
        "POST",
        payload,
        user
      );
      if (!res.data)
        throw new Error(res.message || "Failed to assign homework");

      setToast({ open: true, msg: "Homework assigned!" });
      setSelected([]);
      setTitle("");
      setDue("");
      setDesc("");
      onClose();
    } catch (e) {
      setToast({
        open: true,
        msg: e instanceof Error ? e.message : "Failed to assign",
        error: true,
      });
    } finally {
      setBusy(false);
    }
  };

  const closeToast = () => setToast({ open: false, msg: "" });

  return (
    <>
      <Dialog
        open={open}
        onClose={busy ? undefined : onClose}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Assign homework</DialogTitle>
        <DialogContent sx={{ pt: 1 }}>
          <Stack spacing={2} sx={{ mt: 1 }}>
            <Autocomplete
              multiple
              options={sections}
              loading={loadingSections}
              value={selected}
              onChange={(_, v) => setSelected(v)}
              getOptionLabel={(o) => o.name ?? o.id}
              renderValue={(value) =>
                value.map((option) => (
                  <Chip key={option.id} label={option.name || option.id} />
                ))
              }
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Sections"
                  placeholder="Select sections..."
                  slotProps={{
                    input: {
                      ...params.InputProps,
                      endAdornment: (
                        <>
                          {loadingSections ? (
                            <CircularProgress color="inherit" size={18} />
                          ) : null}
                          {params.InputProps.endAdornment}
                        </>
                      ),
                    },
                  }}
                />
              )}
            />
            <TextField
              label="Title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
              fullWidth
            />
            <TextField
              label="Due date & time"
              type="datetime-local"
              value={due}
              onChange={(e) => setDue(e.target.value)}
              required
              fullWidth
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
              helperText="Local time; will be stored as UTC"
            />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} disabled={busy}>
            Cancel
          </Button>
          <Button variant="contained" onClick={handleSubmit} disabled={busy}>
            {busy ? "Assigning…" : "Assign"}
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={toast.open}
        autoHideDuration={1700}
        onClose={closeToast}
        message={toast.msg}
        anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
        slotProps={{
          content: {
            sx: (t) => ({
              bgcolor: toast.error
                ? t.palette.error.main
                : t.palette.success.main,
            }),
          },
        }}
      />
    </>
  );
}
