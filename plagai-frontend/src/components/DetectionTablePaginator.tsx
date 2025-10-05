"use client";

import { useEffect, useMemo, useState } from "react";
import {
  Box,
  CircularProgress,
  TablePagination,
  Typography,
} from "@mui/material";
import DetectionTable from "./DetectionTable";
import { Detection } from "@/types/types";
import { API_BASE_URL } from "@/lib/api";
import { useAuth } from "@/providers/AuthProvider";

type Props = {
  sectionID: number | string;
  homeworkID: number | string;
  pageSize?: number; // default 10 seemed reasonable, change if needed ~brcrt
  onView: (d: Detection) => void;
};

export default function DetectionTablePaginated({
  sectionID,
  homeworkID,
  pageSize = 10,
  onView,
}: Props) {
  const { user } = useAuth();
  const [rows, setRows] = useState<Detection[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(pageSize);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!sectionID || !homeworkID) return;
    let cancelled = false;
    const run = async () => {
      setLoading(true);
      setError(null);
      try {
        const qp = new URLSearchParams({
          section: String(sectionID),
          homework: String(homeworkID),
          page: String(page + 1),
          limit: String(limit),
        });
        // By the way now that I'm looking back at this, I wonder why I didn't use fetchFromAPI here ~brcrt
        // Oh well, too late now I guess
        const res = await fetch(`${API_BASE_URL}/detections?${qp.toString()}`, {
          method: "GET",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
            ...(user?.token ? { Authorization: `Bearer ${user.token}` } : {}),
          },
        });

        if (!res.ok) {
          const msg = await res.text();
          throw new Error(msg || "Failed to load detections");
        }

        const totalStr = res.headers.get("X-Total-Count") ?? "0";
        const pageStr = res.headers.get("X-Page") ?? String(page + 1);
        const limitStr = res.headers.get("X-Limit") ?? String(limit);

        const body = await res.json();
        if (cancelled) return;

        setRows(body?.data ?? []);
        setTotal(parseInt(totalStr) || 0);

        const srvPage = Math.max(1, parseInt(pageStr) || page + 1);
        const srvLimit = Math.max(1, parseInt(limitStr) || limit);
        if (srvLimit !== limit) setLimit(srvLimit);
        if (srvPage - 1 !== page) setPage(srvPage - 1);
      } catch (e) {
        if (cancelled) return;
        setError(e instanceof Error ? e.message : "Failed to load detections");
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    run();
    return () => {
      cancelled = true;
    };
  }, [API_BASE_URL, sectionID, homeworkID, page, limit, user?.token]);

  const rangeLabel = useMemo(() => {
    if (total === 0) return "0–0 of 0";
    const start = page * limit + 1;
    const end = Math.min(total, start + rows.length - 1);
    return `${start}–${end} of ${total}`;
  }, [page, limit, rows.length, total]);

  return (
    <Box>
      {loading && rows.length === 0 ? (
        <Box sx={{ py: 6, textAlign: "center" }}>
          <CircularProgress />
          <Typography variant="body2" sx={{ mt: 2 }}>
            Loading…
          </Typography>
        </Box>
      ) : error ? (
        <Box sx={{ py: 6, textAlign: "center" }}>
          <Typography variant="h6" gutterBottom>
            Error loading detections
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {error}
          </Typography>
        </Box>
      ) : (
        <>
          <DetectionTable detections={rows} onView={onView} />
          <TablePagination
            component="div"
            count={total}
            page={page}
            onPageChange={(_, newPage) => setPage(newPage)}
            rowsPerPage={limit}
            onRowsPerPageChange={(e) => {
              const next = parseInt(e.target.value, 10);
              setLimit(next > 0 ? next : 50);
              setPage(0);
            }}
            labelDisplayedRows={() => rangeLabel}
            rowsPerPageOptions={[10, 25, 50, 100, 200]}
          />
        </>
      )}
    </Box>
  );
}
