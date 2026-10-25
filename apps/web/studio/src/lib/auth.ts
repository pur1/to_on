import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ApiError, authApi, type User } from "./api";

const currentUserKey = ["auth", "me"] as const;

// Resolves to the signed-in user, or null for guests. On first load the
// access token is gone, so the session is restored from the refresh cookie.
export function useCurrentUser() {
  return useQuery({
    queryKey: currentUserKey,
    queryFn: async (): Promise<User | null> => {
      try {
        return await authApi.refresh();
      } catch (err) {
        if (err instanceof ApiError && err.status === 401) return null;
        throw err;
      }
    },
    staleTime: Infinity,
    retry: false,
  });
}

export function useLogin() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: authApi.login,
    onSuccess: (user) => queryClient.setQueryData(currentUserKey, user),
  });
}

export function useLogout() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: authApi.logout,
    onSettled: () => queryClient.setQueryData(currentUserKey, null),
  });
}
