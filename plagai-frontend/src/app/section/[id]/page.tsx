"use client";
import { useAuth } from "@/providers/AuthProvider";
import AuthGuard from "@/components/AuthGuard";
import { useEffect, useState, useMemo } from "react";
import { fetchFromAPI } from "@/lib/api";
import { useParams, useRouter } from "next/navigation";
import Header from "@/components/Header";
import { Container, Stack, Typography, Alert, Chip } from "@mui/material";
import RefreshIcon from "@mui/icons-material/Refresh";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import CardsGrid from "@/components/CardGrid";
import { Homework, Section } from "@/types/types";
import PageHeader from "@/components/PageHeader";

function SectionHeader({
  loading,
  onRefresh,
  name,
}: {
  loading: boolean;
  onRefresh: () => void;
  name?: string;
}) {
  const router = useRouter();
  const { id } = useParams<{ id: string }>();

  return (
    <PageHeader
      breadcrumbs={[
        { label: "Home", onClick: () => router.push("/") },
        { label: name ? `${name}` : `Section ${id}` },
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

export default function SectionPage() {
  const params = useParams();
  const { user } = useAuth();

  const sectionId = useMemo(
    () =>
      typeof params?.id === "string"
        ? params.id
        : Array.isArray(params?.id)
        ? params.id[0]
        : "",
    [params]
  );

  const [details, setDetails] = useState<Section>();
  const [homeworks, setHomeworks] = useState<Homework[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadHomeworks = async () => {
    try {
      setError(null);
      setLoading(true);
      // This should eventually be fetched per-section when backend is properly setup.
      // Something like /detections?sectionid=1 or whatever
      const res = await fetchFromAPI<Homework[]>(
        `/homeworks?section=${sectionId}`,
        "GET",
        null,
        user
      ); // This type nesting is fucking disgusting but whatever ~brtcrt
      if (!res.data) throw new Error("Failed to fetch homeworks");
      setHomeworks(res.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setLoading(false);
    }
  };

  const loadDetails = async () => {
    try {
      if (!loading) setLoading(true);
      const res = await fetchFromAPI<Section>(
        `/section?section=${sectionId}`,
        "GET",
        null,
        user
      );
      if (!res.data) throw new Error("Failed to get section details");
      setDetails(res.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error");
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => {
    if (!sectionId) return;
    loadHomeworks();
    loadDetails();
  }, [sectionId, user]);

  return (
    <AuthGuard>
      <Header />
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Stack spacing={2} sx={{ mb: 2 }}>
          <SectionHeader
            loading={loading}
            onRefresh={loadHomeworks}
            name={details?.name}
          />

          <Stack direction="row" alignItems="center" spacing={1}>
            <Typography variant="h4" fontWeight={700}>
              {details?.name}
            </Typography>
            <Chip label={`Homeworks: ${homeworks.length}`} size="small" />
          </Stack>

          <Typography variant="body2" color="text.secondary">
            Review homeworks for this section. Click on a card to view
            detections for that homework.
          </Typography>
        </Stack>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        <CardsGrid
          loading={loading}
          emptyText="No homeworks for this section"
          items={homeworks.map((hw) => ({
            id: hw.id,
            title: hw.title,
            subtitle: hw.dueDate
              ? `Due ${new Date(hw.dueDate).toLocaleDateString()}`
              : undefined,
            href: `/section/${sectionId}/homework/${hw.id}`,
          }))}
        ></CardsGrid>
      </Container>
    </AuthGuard>
  );
  // I'm actually about to kill myself holy shit frontend is just pain and suffering. ~brtcrt
}
