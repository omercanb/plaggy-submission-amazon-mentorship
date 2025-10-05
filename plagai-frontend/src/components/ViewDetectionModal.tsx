"use client";

import { useMemo, useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  IconButton,
  Tooltip,
  Divider,
  Snackbar,
  useMediaQuery,
  useTheme,
  Typography,
  Switch,
  FormControlLabel,
} from "@mui/material";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import { alpha } from "@mui/material/styles";

interface ViewParagraphModalProps {
  open: boolean;
  content: string; // this is just flagData btw ~brtcrts
  onClose: () => void;
  filePath?: string;
  flagText?: string;
}

type DiffRow = {
  type: "context" | "add" | "del";
  leftNo?: number;
  rightNo?: number;
  left: string;
  right: string;
};

type Hunk = {
  header: string;
  rows: DiffRow[];
};

// I don't know what the fuck is going on here. I have been staring at this for
// the past 2 hours and am about to pass the fuck out please send help. ~brtcrt
function parseUnifiedDiff(diffRaw: string): Hunk[] {
  const diff = (diffRaw || "").replace(/\r\n/g, "\n");
  const lines = diff.split("\n");

  const hunks: Hunk[] = [];
  let i = 0;

  const pushSimpleHunkIfEmpty = () => {
    // For very short diffs missing @@ headers, we shouldn't need this but just in case ~brtcrt
    if (hunks.length === 0) {
      let leftNo = 1;
      let rightNo = 1;
      const rows: DiffRow[] = [];
      for (const l of lines) {
        if (!l) continue;
        if (l.startsWith("+")) {
          // also I had no way of testing if this works since we don't have a diff that has addition
          // for god knows what reason but I shouldn't have to come back here hopefully and well
          // if I do then look towards the end of the file for my honest thoughts ~brtcrt
          rows.push({
            type: "add",
            rightNo: rightNo++,
            left: "",
            right: l.slice(1),
          });
        } else if (l.startsWith("-")) {
          rows.push({
            type: "del",
            leftNo: leftNo++,
            left: l.slice(1),
            right: "",
          });
        } else {
          const t = l.startsWith(" ") ? l.slice(1) : l;
          rows.push({
            type: "context",
            leftNo: leftNo++,
            rightNo: rightNo++,
            left: t,
            right: t,
          });
        }
      }
      if (rows.length > 0) hunks.push({ header: "", rows });
    }
  };

  while (i < lines.length) {
    const line = lines[i];
    const m = line?.startsWith("@@")
      ? line.match(/^@@\s*-(\d+)(?:,(\d+))?\s+\+(\d+)(?:,(\d+))?\s*@@/)
      : null;

    if (!m) {
      i++;
      continue;
    }

    let leftNo = parseInt(m[1], 10) || 1;
    let rightNo = parseInt(m[3], 10) || 1;
    const hunk: Hunk = { header: line, rows: [] };
    i++;

    while (i < lines.length && !lines[i].startsWith("@@")) {
      const l = lines[i];
      if (l?.startsWith("--- ") || l?.startsWith("+++ ")) break;

      if (l?.startsWith("+")) {
        hunk.rows.push({
          type: "add",
          leftNo: undefined,
          rightNo: rightNo++,
          left: "",
          right: l.slice(1),
        });
      } else if (l?.startsWith("-")) {
        hunk.rows.push({
          type: "del",
          leftNo: leftNo++,
          rightNo: undefined,
          left: l.slice(1),
          right: "",
        });
      } else if (l?.startsWith("\\ No newline at end of file")) {
        // chatgpt suggested this but i don't think this is necessary ~brtcrt
      } else {
        const t = l?.startsWith(" ") ? l.slice(1) : l ?? "";
        hunk.rows.push({
          type: "context",
          leftNo: leftNo++,
          rightNo: rightNo++,
          left: t,
          right: t,
        });
      }
      i++;
    }

    hunks.push(hunk);
  }

  if (hunks.length === 0) pushSimpleHunkIfEmpty();
  return hunks;
}

export default function ViewParagraphModal({
  open,
  content,
  onClose,
  filePath,
  flagText,
}: ViewParagraphModalProps) {
  const theme = useTheme();
  const fullScreen = useMediaQuery(theme.breakpoints.down("sm"));

  const [wrap, setWrap] = useState(true);
  const [copied, setCopied] = useState(false);

  const hunks = useMemo(() => parseUnifiedDiff(content), [content]);

  const stats = useMemo(() => {
    let adds = 0;
    let dels = 0;
    for (const h of hunks) {
      for (const r of h.rows) {
        if (r.type === "add") adds++;
        else if (r.type === "del") dels++;
      }
    }
    return { adds, dels, chars: (content?.length ?? 0).toLocaleString() };
  }, [hunks, content]);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(content || "");
      setCopied(true);
    } catch (e) {
      console.error(`Error while copying to clipboard: ${e}`);
    }
  };
  // believe it or not this was the most cancer part ~brtcrt
  const bgFor = (t: DiffRow["type"]) => {
    if (t === "add")
      return alpha(
        theme.palette.success.light,
        theme.palette.mode === "dark" ? 0.25 : 0.35
      );
    if (t === "del")
      return alpha(
        theme.palette.error.light,
        theme.palette.mode === "dark" ? 0.25 : 0.35
      );
    return theme.palette.mode === "light"
      ? theme.palette.grey[50]
      : theme.palette.grey[900];
  };
  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="lg"
      fullWidth
      fullScreen={fullScreen}
    >
      <DialogTitle sx={{ pb: 1.5 }}>
        <Box
          display="flex"
          alignItems="center"
          justifyContent="space-between"
          gap={1}
        >
          <Box>
            <Typography variant="h6" fontWeight={700}>
              {filePath ? `Changes in ${filePath}` : "Diff Viewer"}
            </Typography>
            <Typography
              variant="caption"
              color="text.secondary"
              sx={{ display: "block" }}
            >
              {stats.chars} chars • +{stats.adds} / −{stats.dels}
            </Typography>
            {flagText && (
              <Typography variant="caption" color="text.secondary">
                Flag: {flagText}
              </Typography>
            )}
          </Box>

          <Box display="flex" alignItems="center" gap={1}>
            <FormControlLabel
              control={
                <Switch
                  size="small"
                  checked={wrap}
                  onChange={(e) => setWrap(e.target.checked)}
                />
              }
              label="Wrap"
              sx={{ m: 0 }}
            />
            <Tooltip title="Copy raw diff">
              <IconButton onClick={handleCopy} size="small" aria-label="Copy">
                <ContentCopyIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          </Box>
        </Box>
      </DialogTitle>

      <Divider />
      {content == "" || !content ? (
        // In case we have to show an empty diff which I think is fucking retarded
        // but half of the diffs we have in the db are empty ~brtcrt
        <DialogContent>
          <Typography variant="body2" color="text.secondary">
            No patch to show.
          </Typography>
        </DialogContent>
      ) : (
        <DialogContent sx={{ p: { xs: 1.5, sm: 2.5 } }}>
          <Box
            sx={{
              display: "grid",
              gridTemplateColumns: "72px 1fr 72px 1fr",
              gap: 0.5,
              px: 1,
              pb: 1,
            }}
          >
            <Box />
            <Typography variant="body2" color="text.secondary">
              Original
            </Typography>
            <Box />
            <Typography variant="body2" color="text.secondary">
              Updated
            </Typography>
          </Box>

          <Box
            sx={{
              border: (t) => `1px solid ${t.palette.divider}`,
              borderRadius: 2,
              overflow: "hidden",
            }}
          >
            {hunks.map((hunk, hi) => (
              <Box key={hi}>
                {hunk.header && (
                  // This is actually really unintuitive but I don't give a shit and
                  // also have no fucking clue how else I could do the background shit ~brtcrt
                  <Box
                    component="pre"
                    sx={{
                      m: 0,
                      px: 1.5,
                      py: 0.75,
                      bgcolor: (t) =>
                        t.palette.mode === "light"
                          ? t.palette.grey[100]
                          : t.palette.grey[800],
                      borderBottom: (t) => `1px solid ${t.palette.divider}`,
                      fontFamily:
                        'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
                      fontSize: 12,
                    }}
                  >
                    {hunk.header}
                  </Box>
                )}

                {hunk.rows.map((row, ri) => (
                  <Box
                    key={`${hi}-${ri}`}
                    sx={{
                      display: "grid",
                      gridTemplateColumns: "72px 1fr 72px 1fr",
                      alignItems: "stretch",
                      borderTop: (t) => `1px solid ${t.palette.divider}`,
                    }}
                  >
                    <Box
                      sx={{
                        px: 1,
                        py: 0.5,
                        textAlign: "right",
                        fontVariantNumeric: "tabular-nums",
                        color: "text.secondary",
                        bgcolor: bgFor(row.type),
                        borderRight: (t) => `1px solid ${t.palette.divider}`,
                      }}
                    >
                      {row.leftNo ?? ""}
                    </Box>
                    <Box
                      component="pre"
                      sx={{
                        m: 0,
                        px: 1.5,
                        py: 0.5,
                        bgcolor: bgFor(row.type),
                        overflowX: "auto",
                        whiteSpace: wrap ? "pre-wrap" : "pre",
                        wordBreak: wrap ? "break-word" : "normal",
                        fontFamily:
                          'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
                        fontSize: 13,
                      }}
                    >
                      {row.left}
                    </Box>

                    <Box
                      sx={{
                        px: 1,
                        py: 0.5,
                        textAlign: "right",
                        fontVariantNumeric: "tabular-nums",
                        color: "text.secondary",
                        bgcolor: bgFor(row.type),
                        borderRight: (t) => `1px solid ${t.palette.divider}`,
                      }}
                    >
                      {row.rightNo ?? ""}
                    </Box>
                    <Box
                      component="pre"
                      sx={{
                        m: 0,
                        px: 1.5,
                        py: 0.5,
                        bgcolor: bgFor(row.type),
                        overflowX: "auto",
                        whiteSpace: wrap ? "pre-wrap" : "pre",
                        wordBreak: wrap ? "break-word" : "normal",
                        fontFamily:
                          'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
                        fontSize: 13,
                      }}
                    >
                      {row.right}
                    </Box>
                  </Box>
                ))}
              </Box>
            ))}
          </Box>
        </DialogContent>
      )}
      {/* I stg if I need to change how this looks I will blow a hole straight through my head ~brtcrt */}

      <DialogActions sx={{ px: { xs: 2, sm: 3 }, pb: { xs: 2, sm: 3 } }}>
        <Button onClick={onClose} variant="contained">
          Close
        </Button>
      </DialogActions>

      <Snackbar
        open={copied}
        autoHideDuration={1600}
        onClose={() => setCopied(false)}
        message="Copied to clipboard"
        anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
      />
    </Dialog>
  );
}
