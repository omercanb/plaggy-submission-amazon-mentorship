"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import AuthGuard from "@/components/AuthGuard";
import Header from "@/components/Header";
import { fetchFromAPI } from "@/lib/api";
import { Detection, Homework, Section } from "@/types/types";
import {
  Container,
  Box,
  Typography,
  CircularProgress,
  Card,
  CardContent,
  Divider,
  Alert,
  Stack,
  Tooltip,
} from "@mui/material";
import DetectionTable from "@/components/DetectionTable";
import ViewParagraphModal from "@/components/ViewDetectionModal";
import { useAuth } from "@/providers/AuthProvider";
import {
  ArrowBack as ArrowBackIcon,
  Refresh as RefreshIcon,
} from "@mui/icons-material";
import PageHeader from "@/components/PageHeader";
import CalendarTodayIcon from "@mui/icons-material/CalendarToday";
import AccessTimeIcon from "@mui/icons-material/AccessTime";
import EventAvailableIcon from "@mui/icons-material/EventAvailable";
import { Chip, Grid } from "@mui/material";
import DetectionTablePaginated from "@/components/DetectionTablePaginator";

// These are pretty useful, probably move them to some other place to reuse
const fmtDateTime = (d: Date) =>
  new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(d);

const fmtRelative = (d: Date) => {
  const rtf = new Intl.RelativeTimeFormat(undefined, { numeric: "auto" });
  const now = Date.now();
  const diffMs = d.getTime() - now;
  const diffDays = Math.round(diffMs / (1000 * 60 * 60 * 24));
  return rtf.format(diffDays, "day");
};

function HomeworkHeader({
  loading,
  onRefresh,
  sectionName,
  hwName,
}: {
  loading: boolean;
  onRefresh: () => void;
  sectionName?: string;
  hwName?: string;
}) {
  const router = useRouter();
  const { id, hwId } = useParams<{ id: string; hwId: string }>();

  return (
    <PageHeader
      breadcrumbs={[
        { label: "Home", onClick: () => router.push("/") },
        {
          label: sectionName ? `${sectionName}` : `Section ${id}`,
          onClick: () => router.push(`/section/${id}`),
        },
        { label: hwName ? `${hwName}` : `Homework ${hwId}` },
      ]}
      actions={[
        {
          title: "Back",
          icon: <ArrowBackIcon fontSize="small" />,
          onClick: () => router.back(),
        },
        {
          title: "Refresh",
          icon: <RefreshIcon fontSize="small" />,
          onClick: onRefresh,
          disabled: loading,
        },
      ]}
    />
  );
}

export default function HomeworkDetectionsPage() {
  const { id, hwId } = useParams<{ id: string; hwId: string }>();
  const [detections, setDetections] = useState<Detection[]>([]);
  const [loading, setLoading] = useState(true);
  const [viewOpen, setViewOpen] = useState(false);
  const [selected, setSelected] = useState<Detection>({
    content: "",
    createdAt: "",
    createdBy: "",
    diffData: "",
    filePath: "",
    homeworkId: "",
    id: "",
  } as Detection);
  const { user } = useAuth();

  const [sectionDetails, setSectionDetails] = useState<Section>();
  const [hwDetails, setHwDetails] = useState<Homework>();
  const [error, setError] = useState<string | null>(null);

  const loadDetections = async () => {
    try {
      const res = await fetchFromAPI<Detection[]>(
        `/detections?section=${id}&homework=${hwId}`,
        "GET",
        null,
        user
      );
      setDetections(res.data || []);
    } finally {
      setLoading(false);
    }
  };

  const loadSectionDetails = async () => {
    try {
      if (!loading) setLoading(true);
      const res = await fetchFromAPI<Section>(
        `/section?section=${id}`,
        "GET",
        null,
        user
      );
      if (!res.data) throw new Error("Failed to get section details");
      setSectionDetails(res.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setLoading(false);
    }
  };

  const loadHwDetails = async () => {
    try {
      if (!loading) setLoading(true);
      const res = await fetchFromAPI<Homework>(
        `/homework?homework=${hwId}`,
        "GET",
        null,
        user
      );
      if (!res.data) throw new Error("Failed to get section details");
      setHwDetails(res.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadDetections();
    loadHwDetails();
    loadSectionDetails();
  }, [id, hwId]);

  return (
    <AuthGuard>
      <Header />
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <HomeworkHeader
          loading={loading}
          onRefresh={loadDetections}
          hwName={hwDetails?.title}
          sectionName={sectionDetails?.name}
        />
        <Stack direction="row" alignItems="center" spacing={1} sx={{ mb: 2 }}>
          <Typography variant="h4" fontWeight={700}>
            {hwDetails?.title || `Homework ${hwId}`}
          </Typography>
          {sectionDetails?.name && (
            <Chip label={sectionDetails.name} size="small" />
          )}
        </Stack>
        <Card
          variant="outlined"
          sx={(t) => ({
            mb: 3,
            borderRadius: 3,
            overflow: "hidden",
            borderColor: t.palette.divider,
            backgroundColor:
              t.palette.mode === "light"
                ? t.palette.background.paper
                : t.palette.background.default,
          })}
        >
          <CardContent sx={{ p: 2.5 }}>
            <Grid container spacing={2}>
              <Grid size={{ xs: 12, md: 4 }}>
                <Stack direction="row" spacing={1.5} alignItems="flex-start">
                  <CalendarTodayIcon fontSize="small" />
                  <Box>
                    <Typography variant="overline" color="text.secondary">
                      Assigned
                    </Typography>
                    {hwDetails?.assignedAt ? (
                      <Tooltip
                        title={new Date(hwDetails.assignedAt).toISOString()}
                      >
                        <Typography variant="body2">
                          {fmtDateTime(new Date(hwDetails.assignedAt))}
                        </Typography>
                      </Tooltip>
                    ) : (
                      <Typography variant="body2">—</Typography>
                    )}
                  </Box>
                </Stack>
              </Grid>

              <Grid size={{ xs: 12, md: 4 }}>
                <Stack direction="row" spacing={1.5} alignItems="flex-start">
                  <EventAvailableIcon fontSize="small" />
                  <Box>
                    <Typography variant="overline" color="text.secondary">
                      Due
                    </Typography>
                    {hwDetails?.dueDate ? (
                      <Tooltip
                        title={new Date(hwDetails.dueDate).toISOString()}
                      >
                        <Typography variant="body2">
                          {fmtDateTime(new Date(hwDetails.dueDate))}
                        </Typography>
                      </Tooltip>
                    ) : (
                      <Typography variant="body2">—</Typography>
                    )}
                  </Box>
                </Stack>
              </Grid>

              <Grid size={{ xs: 12, md: 4 }}>
                <Stack direction="row" spacing={1.5} alignItems="flex-start">
                  <AccessTimeIcon fontSize="small" />
                  <Box>
                    <Typography variant="overline" color="text.secondary">
                      Status
                    </Typography>
                    {hwDetails?.dueDate ? (
                      (() => {
                        const due = new Date(hwDetails.dueDate);
                        const now = new Date();
                        const overdue = due.getTime() < now.getTime();
                        return (
                          <Stack
                            direction="row"
                            spacing={1}
                            alignItems="center"
                          >
                            <Chip
                              size="small"
                              label={overdue ? "Overdue" : "Upcoming"}
                              color={overdue ? "error" : "success"}
                              variant="outlined"
                            />
                            <Typography variant="body2" color="text.secondary">
                              {fmtRelative(due)}
                            </Typography>
                          </Stack>
                        );
                      })()
                    ) : (
                      <Typography variant="body2">—</Typography>
                    )}
                  </Box>
                </Stack>
              </Grid>
            </Grid>
          </CardContent>
        </Card>

        <Card variant="outlined" sx={{ borderRadius: 3, overflow: "hidden" }}>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          <CardContent sx={{ p: 0 }}>
            <Box sx={{ px: 2.5, py: 2 }}>
              <Typography variant="h6">Detections</Typography>
            </Box>
            <Divider />
            {loading ? (
              <Box sx={{ display: "flex", justifyContent: "center", py: 6 }}>
                <CircularProgress />
              </Box>
            ) : (
              <Box sx={{ p: 2 }}>
                <DetectionTablePaginated
                  sectionID={Number(id)}
                  homeworkID={Number(hwId)}
                  onView={(detection) => {
                    setSelected(detection);
                    setViewOpen(true);
                  }}
                />
              </Box>
            )}
          </CardContent>
        </Card>

        <ViewParagraphModal
          open={viewOpen}
          content={selected.diffData}
          onClose={() => setViewOpen(false)}
          filePath={selected.filePath}
          flagText={selected.content}
        />
      </Container>
    </AuthGuard>
  );
}
