import { useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/axios';
import { AuthResponse, LoginInput, RegisterInput, ForgotPasswordInput, User } from '../types';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';

export const useLogin = () => {
  const navigate = useNavigate();
  return useMutation({
    mutationFn: async (data: LoginInput) => {
      const response = await api.post<AuthResponse>('/auth/login', data);
      return response.data;
    },
    onSuccess: (data) => {
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
      toast.success('Login successful');
      navigate('/dashboard');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Login failed');
    },
  });
};

export const useRegister = () => {
  const navigate = useNavigate();
  return useMutation({
    mutationFn: async (data: RegisterInput) => {
      const response = await api.post<AuthResponse>('/auth/register', data);
      return response.data;
    },
    onSuccess: () => {
      toast.success('Registration successful! Please login.');
      navigate('/login');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Registration failed');
    },
  });
};

export const useForgotPassword = () => {
  return useMutation({
    mutationFn: async (data: ForgotPasswordInput) => {
      const response = await api.post('/auth/forgot-password', data);
      return response.data;
    },
    onSuccess: (data: any) => {
      toast.success(data.message || 'If email exists, reset link sent.');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Request failed');
    },
  });
};

export const useLogout = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  
  return () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    queryClient.clear();
    navigate('/login');
    toast.success('Logged out successfully');
  };
};

export const useUser = () => {
  return {
    user: JSON.parse(localStorage.getItem('user') || 'null') as User | null,
    isAuthenticated: !!localStorage.getItem('token'),
  };
};
