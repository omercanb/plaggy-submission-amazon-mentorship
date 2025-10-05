"use client";

import { Detection } from "@/types/types";
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Tooltip,
  Chip,
  Box,
  Typography,
} from "@mui/material";
import { Visibility } from "@mui/icons-material";
import { useMemo } from "react";

interface DetectionTableProps {
  detections: Detection[];
  onView: (detection: Detection) => void;
}
/**
 * @author brtcrt
 * @description A simple table for displaying suspicious code snippets.
 * @param Detections
 * @param onView
 * @returns
 */
export default function DetectionTable({
  detections,
  onView,
}: DetectionTableProps) {
  const rows = useMemo(
    () =>
      [...(detections ?? [])].sort(
        (a, b) =>
          Date.parse(b.createdAt as unknown as string) -
          Date.parse(a.createdAt as unknown as string)
      ),
    [detections]
  );

  const hasRows = rows.length > 0;

  const formatDisplayDate = (iso: string | Date) =>
    new Intl.DateTimeFormat(undefined, {
      dateStyle: "medium",
      timeStyle: "short",
    }).format(new Date(iso));

  return (
    <Paper
      elevation={3}
      sx={{
        borderRadius: 3,
        overflow: "hidden",
      }}
    >
      {!hasRows ? (
        <Box sx={{ py: 6, textAlign: "center" }}>
          <Typography variant="h6" gutterBottom>
            No detections yet
          </Typography>
          <Typography variant="body2" color="text.secondary">
            When detections are available, they’ll show up here.
          </Typography>
        </Box>
      ) : (
        <TableContainer>
          <Table stickyHeader aria-label="detections table" size="small">
            <TableHead>
              <TableRow>
                <TableCell sx={{ fontWeight: 700, width: 160 }}>
                  Student
                </TableCell>
                <TableCell sx={{ fontWeight: 700 }}>
                  Content Preview
                </TableCell>
                <TableCell sx={{ fontWeight: 700, width: 200 }}>
                  Created
                </TableCell>
                <TableCell sx={{ fontWeight: 700, width: 60 }} align="center">
                  Severity
                </TableCell>
                <TableCell sx={{ fontWeight: 700, width: 100, textAlign: "right" }}>
                  Actions
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {rows.map((d) => {
                const fullDate = new Date(
                  d.createdAt as unknown as string
                ).toISOString();
                const displayDate = formatDisplayDate(
                  d.createdAt as unknown as string
                );
                const preview = d.content?.trim() ?? "";
                const truncated =
                  preview.length > 120
                    ? `${preview.slice(0, 120).replace(/\s+\S*$/, "")}…`
                    : preview;

                // Map severity to color
                const severityColor = (() => {
                  switch (d.severity) {
                    case 1:
                      return "default";
                    case 2:
                      return "warning";
                    case 3:
                      return "error";
                    default:
                      return "default";
                  }
                })();

                return (
                  <TableRow
                    key={d.id}
                    hover
                    sx={{
                      "&:nth-of-type(odd)": {
                        backgroundColor: (t) => t.palette.action.hover,
                      },
                    }}
                  >
                    <TableCell>
                      <Chip
                        label={d.createdBy || "Unknown"}
                        size="small"
                        variant="outlined"
                        sx={{ borderRadius: 2 }}
                      />
                    </TableCell>

                    <TableCell sx={{ maxWidth: 0 }}>
                      <Tooltip title={preview} placement="top" enterDelay={300}>
                        <Typography
                          variant="body2"
                          sx={{
                            whiteSpace: "nowrap",
                            overflow: "hidden",
                            textOverflow: "ellipsis",
                          }}
                        >
                          {truncated || "—"}
                        </Typography>
                      </Tooltip>
                    </TableCell>

                    <TableCell>
                      <Tooltip
                        title={fullDate}
                        placement="top"
                        enterDelay={300}
                      >
                        <Typography variant="body2">{displayDate}</Typography>
                      </Tooltip>
                    </TableCell>

                    <TableCell align="center">
                      <Chip
                        label={d.severity}
                        color={severityColor}
                        size="small"
                        sx={{ fontWeight: 600 }}
                      />
                    </TableCell>

                    <TableCell align="right">
                      <Tooltip title="View full content">
                        <IconButton
                          onClick={() => onView(d)}
                          size="small"
                          aria-label="view content"
                        >
                          <Visibility fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </TableContainer>
      )}
    </Paper>
  );
}
