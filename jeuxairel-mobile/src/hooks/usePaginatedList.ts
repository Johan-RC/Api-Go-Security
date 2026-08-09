import { useCallback, useEffect, useRef, useState } from 'react';
import { getErrorMessage } from '@/api/errors';
import type { ListEnvelope } from '@/types/api';

interface PaginatedOptions {
  pageSize?: number;
  initialPage?: number;
  enabled?: boolean;
}

export function usePaginatedList<T>(
  fetcher: (page: number, pageSize: number) => Promise<ListEnvelope<T>>,
  { pageSize = 20, initialPage = 1, enabled = true }: PaginatedOptions = {},
) {
  const [items, setItems] = useState<T[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(initialPage);
  const [loading, setLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);

  const mounted = useRef(true);
  useEffect(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
    };
  }, []);

  const loadPage = useCallback(
    async (targetPage: number, append: boolean) => {
      if (!enabled) return;
      setError(null);
      try {
        const result = await fetcher(targetPage, pageSize);
        if (!mounted.current) return;
        setItems((prev) => (append ? [...prev, ...result.items] : result.items));
        setTotal(result.total);
        setPage(targetPage);
        setHasMore(result.items.length >= pageSize && targetPage * pageSize < result.total);
      } catch (err) {
        if (!mounted.current) return;
        setError(getErrorMessage(err));
      }
    },
    [enabled, fetcher, pageSize],
  );

  useEffect(() => {
    setLoading(true);
    loadPage(initialPage, false).finally(() => mounted.current && setLoading(false));
  }, [loadPage, initialPage]);

  const refresh = useCallback(async () => {
    setRefreshing(true);
    await loadPage(initialPage, false);
    setRefreshing(false);
  }, [loadPage, initialPage]);

  const loadMore = useCallback(() => {
    if (!hasMore || loading || refreshing) return;
    setLoading(true);
    loadPage(page + 1, true).finally(() => mounted.current && setLoading(false));
  }, [hasMore, loading, refreshing, loadPage, page]);

  return { items, total, page, loading, refreshing, error, hasMore, refresh, loadMore };
}