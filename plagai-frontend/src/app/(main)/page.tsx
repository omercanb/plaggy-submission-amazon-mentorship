/**
 * OK this is for whichever sad fucker has to work on this part with me.
 * We are using Next.js folder/file based routing, which means that
 * file naming is very sensative and prone to exploding if you misspell
 * and routing is directly based on the folders, unless you use (), which
 * just groups them without affecting the route.
 *
 * page.tsx is where you group components that eventually get served to the user
 * so for example this page.tsx is for localhost:3000 directly, however /(auth)/login/page.tsx
 * is for localhost:3000/login. We will modify this most of the time, unless we are doing some
 * wonky shit that is required by most pages and should be checked multiple times, like
 * auth, for example.
 *
 * Those are grouped under providers (similar to context in a regular react project) and are then
 * wrapped around the whole entire project structe in layout.tsx, which is like the entry point
 * for the application.
 *
 * ~brtcrt
 *
 */

"use client";

import { useEffect, useState, useMemo } from "react";
import AuthGuard from "@/components/AuthGuard";
import Header from "@/components/Header";
import { fetchFromAPI } from "@/lib/api";
import { Section } from "@/types/types";
import { useRouter } from "next/navigation";
import { useAuth } from "@/providers/AuthProvider";
import {
  Container,
  Typography,
  Stack,
  Chip,
  Alert,
  IconButton,
} from "@mui/material";
import RefreshIcon from "@mui/icons-material/Refresh";
import CardsGrid from "@/components/CardGrid";

export default function HomePage() {
  const { user } = useAuth();

  const [sections, setSections] = useState<Section[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadSections = async () => {
    try {
      setError(null);
      setLoading(true);
      const res = await fetchFromAPI<Section[]>("/sections", "GET", null, user);
      setSections(res.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load sections");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!user) return;
    loadSections();
  }, [user]);

  const countLabel = useMemo(() => {
    const n = sections.length;
    if (loading) return "Loading…";
    if (n === 0) return "No sections";
    return `${n} section${n > 1 ? "s" : ""}`;
  }, [sections.length, loading]);

  return (
    <AuthGuard>
      <Header />
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Stack
          direction={{ xs: "column", sm: "row" }}
          alignItems={{ xs: "flex-start", sm: "center" }}
          justifyContent="space-between"
          spacing={1.5}
          sx={{ mb: 3 }}
        >
          <Stack direction="row" spacing={1.5} alignItems="center">
            <Typography variant="h4" fontWeight={700}>
              Sections
            </Typography>
            <Chip label={countLabel} size="small" />
          </Stack>
          <Stack direction="row" spacing={1}>
            <IconButton onClick={loadSections} disabled={loading} size="small">
              <RefreshIcon fontSize="small" />
            </IconButton>
          </Stack>
        </Stack>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        <CardsGrid
          loading={loading}
          items={sections.map((s) => ({
            id: s.id,
            title: s.name,
            subtitle: `ID: ${s.id}`,
            href: `/section/${s.id}`,
          }))}
          emptyText="No sections to show yet."
        ></CardsGrid>
      </Container>
    </AuthGuard>
  );
  // This is actually kinda pretty idk. Still incredibly painful though. ~brtcrt
  // Idk about the color palette tho. I think I might have punched gpt way too hard
  // on the head, but I think I managed to finally put some sense of visual harmony
  // it that fucker head.
}
