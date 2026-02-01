import { useQuery, keepPreviousData } from '@tanstack/react-query';
import api from '@/lib/axios';

type User = {
  id: number;
  name: string;
  email: string;
  created_at: string;
};

type UserListResponse = {
  data: User[];
  meta: {
    total: number;
    page: number;
    limit: number;
  };
};

type PaginationParams = {
  page: number;
  limit: number;
  search: string;
};

export const useUsers = (params: PaginationParams) => {
  return useQuery({
    queryKey: ['users', params],
    queryFn: async () => {
      const response = await api.get<UserListResponse>('/users', { params });
      return response.data;
    },
    placeholderData: keepPreviousData,
  });
};

